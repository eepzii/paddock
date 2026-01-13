package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eepzii/paddock/internal/app"
	"github.com/eepzii/paddock/internal/browser"
	"github.com/eepzii/paddock/internal/storage"
)

func main() {
	startTime := time.Now()

	flags := flag.NewFlagSet("paddock", flag.ExitOnError)

	browserPathFlag := flags.String("path", "", "sets a custom browser path")
	emailFlag := flags.String("email", "", "sets email")
	freshnessFlag := flags.Duration("freshness", time.Hour, "checks the tokens validity in x future (max freshness is 96h)")
	headlessFlag := flags.Bool("headless", false, "disables visual ui")
	logoutFlag := flags.Bool("logout", false, "logs out->  deletes subscriptionToken-> deletes browser-profile")
	forceFlag := flags.Bool("force", false, "forces login")
	flags.Parse(os.Args[1:])

	log.SetFlags(0)

	if *forceFlag && *logoutFlag {
		app.PrintFatal(errors.New("cannot use --force with --logout"), startTime)
	}

	fileManager, err := storage.New()
	if err != nil {
		app.PrintFatal(err, startTime)
	}
	config, err := fileManager.LoadConfig()
	if err != nil {
		app.PrintFatal(fmt.Errorf("failed to load config: %v", err), startTime)
	}

	bConfig := browser.Config{
		FileManager:       fileManager,
		CustomBrowserPath: *browserPathFlag,
		Headless:          *headlessFlag,
	}
	proxyHttpAddress := os.Getenv("PROXY_HOST")
	if proxyHttpAddress != "" {
		bConfig.Proxy = browser.Proxy{
			Address:  proxyHttpAddress,
			User:     os.Getenv("PROXY_USER"),
			Password: os.Getenv("PROXY_PASS"),
		}
	}

	b, err := browser.New(bConfig)
	if err != nil {
		app.PrintFatal(err, startTime)
	}

	if !*forceFlag && *logoutFlag {
		if err := app.PerformLogout(&b, *emailFlag, config); err != nil {
			app.PrintFatal(err, startTime)
		}
		if err := fileManager.Reset(); err != nil {
			app.PrintFatal(err, startTime)
		}
		app.PrintSuccess("", startTime)
		return
	}

	if err := app.PerformEmailCheck(*emailFlag); err != nil {
		app.PrintFatal(err, startTime)
	}

	if !*forceFlag && app.IsStoredTokenFresh(*emailFlag, config, *freshnessFlag) {
		app.PrintSuccess(config.SubscriptionToken, startTime)
		return
	} else if config.SubscriptionToken != "" && *emailFlag == config.Email {
		if err := app.PerformLogout(&b, *emailFlag, config); err != nil {
			app.PrintFatal(fmt.Errorf("token expired tried to log out: %v", err), startTime)
		}
		if err := fileManager.Reset(); err != nil {
			app.PrintFatal(err, startTime)
		}
	}

	if !*forceFlag && *emailFlag != config.Email && config.Email != "" {
		app.PrintFatal(errors.New("another account is logged in"), startTime)
	}

	password := os.Getenv("PASSWORD")
	if strings.TrimSpace(password) == "" {
		app.PrintFatal(errors.New("set password"), startTime)
	}

	if err := fileManager.Reset(); err != nil {
		app.PrintFatal(err, startTime)
	}

	config, err = app.PerformLogin(&b, *emailFlag, password)
	if err != nil {
		app.PrintFatal(err, startTime)
	}
	if err := fileManager.SaveConfig(config); err != nil {
		app.PrintFatal(err, startTime)
	}

	app.PrintSuccess(config.SubscriptionToken, startTime)
}
