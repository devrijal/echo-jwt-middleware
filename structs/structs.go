package structs

type JWK struct {
	Alg     string   `json:"alg"`
	KeyType string   `json:"kty"`
	KeyID   string   `json:"kid"`
	Use     string   `json:"use"`
	N       string   `json:"n"`
	E       string   `json:"e"`
	X5C     []string `json:"x5c"`
	X5T     string   `json:"x5t"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type OpenIDConfig struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	IntrospectionEndpoint string `json:"introspection_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
	EndSessionEndpoint    string `json:"end_session_endpoint"`
	JWKSUri               string `json:"jwks_uri"`
	RegistrationEndpoint  string `json:"registration_endpoint"`
}
