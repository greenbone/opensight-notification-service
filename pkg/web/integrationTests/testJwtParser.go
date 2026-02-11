// SPDX-FileCopyrightText: 2025 Greenbone AG
//
// SPDX-License-Identifier: GPL-3.0-or-later

package integrationTests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/greenbone/keycloak-client-golang/auth"
)

func NewTestJwtParser(t *testing.T) func(
	ctx context.Context,
	authorizationHeader string,
	originHeader string,
) (auth.UserContext, error) {
	return JwtParserForTesting
}

func JwtParserForTesting(
	ctx context.Context,
	authorizationHeader string,
	originHeader string,
) (auth.UserContext, error) {
	token, err := parseAuthorizationHeader(authorizationHeader)
	if err != nil {
		return auth.UserContext{}, fmt.Errorf("couldn't parse authorization header: %w", err)
	}

	userCtx, err := parseJWT(ctx, token)
	if err != nil {
		return auth.UserContext{}, fmt.Errorf("couldn't parse token: %w", err)
	}

	if originHeader != "" {
		correctOrigin := false
		for _, origin := range userCtx.AllowedOrigins {
			if originHeader == origin {
				correctOrigin = true
				break
			}
		}
		if !correctOrigin {
			return auth.UserContext{}, fmt.Errorf("not allowed origin: %s", originHeader)
		}
	}

	return userCtx, nil
}

func parseJWT(ctx context.Context, token string) (auth.UserContext, error) {
	type customClaims struct {
		jwt.RegisteredClaims
		UserId         string   `json:"sub"`
		Email          string   `json:"email"`
		UserName       string   `json:"preferred_username"`
		Roles          []string `json:"roles"`
		Groups         []string `json:"groups"`
		AllowedOrigins []string `json:"allowed-origins"`
	}

	// token is not validated against the keycloak for testing
	jwtToken, _, err := jwt.NewParser().ParseUnverified(token, &customClaims{})
	if err != nil {
		return auth.UserContext{}, fmt.Errorf("parsing of token failed: %w", err)
	}
	claims := jwtToken.Claims.(*customClaims)

	return auth.UserContext{
		Realm:          "testRealm",
		UserID:         claims.UserId,
		UserName:       claims.UserName,
		EmailAddress:   claims.Email,
		Roles:          claims.Roles,
		Groups:         claims.Groups,
		AllowedOrigins: claims.AllowedOrigins,
	}, nil
}

func parseAuthorizationHeader(authorizationHeader string) (string, error) {
	fields := strings.Fields(authorizationHeader)
	if len(fields) != 2 {
		return "", fmt.Errorf("header contains invalid number of fields: %d", len(fields))
	}
	if strings.ToLower(fields[0]) != "bearer" {
		return "", fmt.Errorf("header contains invalid token type: %q", fields[0])
	}
	return fields[1], nil
}
