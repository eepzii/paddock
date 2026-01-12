package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/eepzii/paddock/internal/browser"
	"github.com/eepzii/paddock/internal/f1site"
	"github.com/eepzii/paddock/internal/storage"
	"github.com/eepzii/paddock/internal/validate"
)

func IsStoredTokenFresh(email string, config storage.Config, duration time.Duration) bool {
	var isFresh = config.TokenExpiration.After(time.Now().Add(duration))
	if email == config.Email && config.SubscriptionToken != "" && isFresh {
		return true
	}
	return false
}

func PerformEmailCheck(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("no email set")
	}

	re := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if !re.MatchString(email) {
		return fmt.Errorf("invalid email structure")
	}

	return nil
}

func PerformLogin(b *browser.Browser, email, password string) (storage.Config, error) {
	loginFunc, loginResultChan := f1site.Login(email, password)
	err := b.Run(loginFunc)
	if err != nil {
		return storage.Config{}, fmt.Errorf("failed to log in: %v", err)
	}

	result := <-loginResultChan
	token, err := validate.JWTToken(result.Response.Body())
	if err != nil {
		return storage.Config{}, fmt.Errorf("failed to verify jwt token: %v", err)
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return storage.Config{}, fmt.Errorf("token exceeds expiration date: %v", err)
	}

	var config = storage.Config{
		Email:             email,
		SubscriptionToken: token.Raw,
		TokenExpiration:   expirationTime.Time,
	}

	return config, nil
}

func PerformLogout(b *browser.Browser, email string, config storage.Config) error {

	if config.SubscriptionToken != "" && email == config.Email {
		if err := b.Run(f1site.Logout); err != nil {
			return fmt.Errorf("failed to log out: %v", err)
		}
	} else if email != config.Email && config.Email != "" {
		return fmt.Errorf("another account is logged in")
	} else {
		return fmt.Errorf("log in first")
	}

	return nil
}
