package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TodoList(c *gin.Context) {
	list := []string{"a", "b", "c"}

	c.JSON(http.StatusOK, list)
}
