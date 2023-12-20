package site

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	t := gin.H{
		"title":     "Todo list",
		"home_page": "Home page",
	}

	c.HTML(http.StatusOK, "home", gin.H{
		"t": t,
	})
}
