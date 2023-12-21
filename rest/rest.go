package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func extractBody(c *gin.Context, data *map[string]interface{}) error {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(value), &data)
	return nil
}

func fullUrl(c *gin.Context) string {
	return c.Request.Host + c.Request.URL.String()
}

func RootIndex(c *gin.Context) {
	c.String(http.StatusOK, "it works")
}
