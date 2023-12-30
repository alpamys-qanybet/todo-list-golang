package api

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"todo/internal/config"
	"todo/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func authorizeToken(c *gin.Context) (uint16, error) { // returns authorized user id
	bearerToken := c.Request.Header.Get("Authorization")

	if bearerToken == "" {
		return 0, errors.New("no token")
	}
	tokenString := strings.Split(bearerToken, " ")[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret()), nil
	})

	if err != nil {
		return 0, err
	}

	return uint16(claims["_content"].(float64)), nil
}

func UserLogin(c *gin.Context) {
	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	login := bodyData["login"].(string)
	password := bodyData["password"].(string)

	if strings.Trim(login, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnprocessableEntity, "login_incorrect_credentials")
		return
	}

	accessToken, err := controller.Authenticate(login, password)
	if err != nil {
		errorMsg := err.Error()

		if errorMsg == "login_incorrect_credentials" {
			c.JSON(http.StatusUnprocessableEntity, errorMsg)
			return
		} else { // else unknown error
			c.JSON(http.StatusUnprocessableEntity, "unknown_error_occurred")
			return
		}
	}

	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfY29udGVudCI6MSwiX3RpbWUiOjE3MDMxNzU2NTM5OTEsIl90b2tlbl9pZCI6IjZlMzY3ODZiLTk2OTItNDllYy1iMTMzLTEzMjRhOGU0YTgwOCJ9.bfKr3WXWD_rwnUVYaIJL-EC4vaeWDMaZaSL_VFaU55w
	// michael
	// jordan

	c.JSON(http.StatusOK, accessToken)
}
