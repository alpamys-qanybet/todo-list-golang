package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootIndex(c *gin.Context) {
	c.String(http.StatusOK, "it works")
}
