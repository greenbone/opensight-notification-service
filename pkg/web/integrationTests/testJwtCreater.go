// SPDX-FileCopyrightText: 2025 Greenbone AG
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package integrationTests

import (
	rand "crypto/rand"
	rsa "crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

const (
	publicKeyID  = "OMTg5TWEm1TZeqeb2zuJJFX1ZxOwDs_IfPIgJ0uIFU0"
	publicKeyALG = "RS256"

	realmId   = "opensight"
	groupId   = "22222222-2222-2222-2222-222222222222"
	publicUrl = "http://localhost/auth"
	origin    = "http://localhost:3000"
)

func getToken(claims jwt.MapClaims, privateKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(publicKeyALG), claims)
	token.Header["kid"] = publicKeyID

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func newPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func CreateJwtTokenWithRole(role string) string {
	return getToken(jwt.MapClaims{
		"iss":                publicUrl + "/realms/" + realmId,
		"sub":                "1927ed8a-3f1f-4846-8433-db290ea5ff90",
		"email":              "user@host.local",
		"preferred_username": "user",
		"roles":              []string{"offline_access", "uma_authorization", role, "default-roles-" + realmId},
		"groups":             []string{groupId},
		"allowed-origins":    []string{origin},
	}, newPrivateKey())
}
