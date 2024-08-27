package common

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/devrijal/jwt-middleware/structs"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetMatchedKey(token *jwt.Token, keys []structs.JWK) (publicKey interface{}, err error) {

	kidInter, ok := token.Header["kid"]

	if !ok {
		return nil, fmt.Errorf("token doesn't have valid kid")
	}

	kid, ok := kidInter.(string)

	if !ok {
		return nil, fmt.Errorf("could not convert kid to string")
	}

	var keyType string

	for _, JWK := range keys {
		if JWK.KeyID == kid && JWK.Use == "sig" {

			n, err := base64.RawURLEncoding.DecodeString(JWK.N)

			if err != nil {
				log.Printf("Error while decode N of key %s: %v", kid, err)
				return nil, err
			}

			e, err := base64.RawURLEncoding.DecodeString(JWK.E)

			if err != nil {
				log.Printf("Error while decode E of key %s: %v", kid, err)
				return nil, err
			}

			keyType = JWK.KeyType

			switch JWK.KeyType {
			case "RSA":
				var pubKey *rsa.PublicKey

				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				pubKey = &rsa.PublicKey{
					N: new(big.Int).SetBytes(n),
					E: int(new(big.Int).SetBytes(e).Int64()),
				}

				publicKey = pubKey
			}
		}
	}

	if publicKey == nil {
		msg := fmt.Sprintf("no suitable %s public key found", keyType)

		return nil, fmt.Errorf(msg)
	}

	return publicKey, nil
}

func GetPublicKeys() (jwks structs.JWKS, err error) {

	client := resty.New()

	openIDConfig, err := GetOpenIDConfig()

	if err != nil {
		return jwks, err
	}

	resp, err := client.R().Get(openIDConfig.JWKSUri)

	if err != nil {
		return jwks, err
	}

	err = json.Unmarshal(resp.Body(), &jwks)

	return jwks, err
}

func GetOpenIDConfig() (config *structs.OpenIDConfig, err error) {

	provider_endpoint := os.Getenv("OPENID_PROVIDER_ENDPOINT")

	if provider_endpoint == "" {
		return config, fmt.Errorf("OpenID provider endpoint not found")
	}

	client := resty.New()

	discovery_endpoint := os.Getenv("OPENID_PROVIDER_DISCOVERY_ENDPOINT")

	if discovery_endpoint == "" {
		return config, fmt.Errorf("OpenID discovery endpoint not found")
	}

	var resp *resty.Response

	resp, err = client.R().Get(discovery_endpoint)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(resp.Body(), &config)

	if err != nil {
		return config, err
	}

	return config, err
}

func GetToken(c echo.Context) *jwt.Token {
	return c.Get("user").(*jwt.Token)
}
