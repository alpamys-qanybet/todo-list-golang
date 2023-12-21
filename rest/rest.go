package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"todo/config"
)

const (
	appSecretName      = "app_secret"
	appSecretIncorrect = "INCORRECT_SECRET"
)

var appSecret string

func SetAppSecret(secret string) {
	appSecret = secret
}

func AppSecret() string {
	return appSecret
}

func extractBody(c *gin.Context, data *map[string]interface{}) error {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(value), &data)
	return nil
}

func appSecretIsValid(c *gin.Context) bool {
	secret := c.Query(appSecretName)

	if secret != AppSecret() {
		if config.DebugLog() {
			log.Printf("app secret is incorrect '%s', must be '%s'", secret, AppSecret())
		}

		c.JSON(http.StatusOK, gin.H{
			"data": appSecretIncorrect,
		})

		return false
	}

	return true
}

func fullUrl(c *gin.Context) string {
	return c.Request.Host + c.Request.URL.String()
}

func RootIndex(c *gin.Context) {
	c.String(http.StatusOK, "it works")
}
