package controller

import (
	"time"

	"todo/internal/config"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func Authenticate(login string, password string) (string, error) {
	id := 1

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	payload := jwt.MapClaims{}
	payload["_time"] = time.Now().UnixMilli()
	payload["_content"] = id
	payload["_token_id"] = tokenId

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(config.JwtSecret()))
	if err != nil {
		return "", err
	}

	return token, nil
}
