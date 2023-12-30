package config

var jwtSecret string

func SetJwtSecret(secret string) {
	jwtSecret = secret
}

func JwtSecret() string {
	return jwtSecret
}
