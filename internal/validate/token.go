package validate

import (
	"context"
	"encoding/json"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTToken(body string) (*jwt.Token, error) {
	var loginRes = LoginResponse{}
	if err := json.Unmarshal([]byte(body), &loginRes); err != nil {
		return nil, err
	}
	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{JWKS_URL})
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(loginRes.Data.SubscriptionToken, jwks.Keyfunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}
