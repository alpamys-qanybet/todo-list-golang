package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"todo/rest"
	"todo/site"
)

func main() {
	log.Println("Todo list App init")

	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*")
	r.GET("/", site.Home)
	r.GET("/rest", rest.RootIndex)
	r.GET("/rest/todo", rest.TodoList)
	r.Run() // listen and serve on "localhost:8080"

	log.Println("Todo list App started SUCCESSFULL")
}
