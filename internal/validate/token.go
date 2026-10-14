package validate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTToken(body string) (*jwt.Token, error) {
	var loginRes = LoginResponse{}
	if err := json.Unmarshal([]byte(body), &loginRes); err != nil {
		return nil, err
	}

	idformula, err := keyfunc.NewDefaultCtx(context.Background(), []string{ID_FORMULA_JWKS})
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(loginRes.AccessToken, claims, idformula.Keyfunc); err != nil {
		return nil, err
	}

	claim, ok := claims[CLAIM].(string)
	if !ok {
		return nil, fmt.Errorf("failed to find subscription token in claims")
	}

	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{API_FORMULA_JWKS})
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(claim, jwks.Keyfunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}
