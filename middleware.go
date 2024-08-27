package echojwtmiddleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Generic middleware
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

		token, err := jwt.Parse(tokenString, GetKey)

		if err != nil || !token.Valid {

			http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Echo jwt.Keyfunc implementation.
var GetKey jwt.Keyfunc = func(token *jwt.Token) (pubKey interface{}, err error) {
	jwks, err := GetPublicKeys()

	if err != nil {
		return nil, err
	}

	pubKey, err = GetMatchedKey(token, jwks.Keys)

	if err != nil {
		return pubKey, err
	}

	return pubKey, nil
}

// Echo middleware.Skipper implementation. Accept paths to skip arguments
func Skipper(pathToSkips []string) func(c echo.Context) bool {

	return func(c echo.Context) bool {

		path := c.Request().URL.Path

		for _, pattern := range pathToSkips {
			re := regexp.MustCompile(pattern)

			isMatch := re.MatchString(path)

			if isMatch {
				return true
			}
		}

		return false
	}
}
