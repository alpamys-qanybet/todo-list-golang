package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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

func RootIndex(c *gin.Context) {
	c.String(http.StatusOK, "it works")
}
