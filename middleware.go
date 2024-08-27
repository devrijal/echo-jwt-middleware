package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/devrijal/jwt-middleware/common"
	"github.com/golang-jwt/jwt/v5"
)

func KeycloakMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, GetMatchedKey)

		if err != nil || !token.Valid {

			http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var GetMatchedKey jwt.Keyfunc = func(token *jwt.Token) (pubKey interface{}, err error) {
	jwks, err := common.GetPublicKeys()

	if err != nil {
		return nil, err
	}

	pubKey, err = common.GetMatchedKey(token, jwks.Keys)

	if err != nil {
		return pubKey, err
	}

	return pubKey, nil
}
