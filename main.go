package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"todo/db"
	"todo/rest"
	"todo/site"
)

var databaseUrl string

func readEnvVariables() (serverHost, serverPort string) {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
	}

	serverHost = os.Getenv("SERVER_HOST")

	serverPort = os.Getenv("SERVER_PORT")
	if "" == serverPort {
		serverPort = "9292"
	}

	databaseUrl = os.Getenv("DATABASE_URL")
	if "" == databaseUrl {
		databaseUrl = "postgresql://postgres:postgres@localhost:5432/todo"
	}

	log.Println("environtment variables are read")
	return
}

func main() {
	log.Println("Todo list App init")

	serverHost, serverPort := readEnvVariables()

	fmt.Println("databaseUrl")
	fmt.Println(databaseUrl)

	dbpool, err := db.Connect(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to postgres database: %v\n", err)
	}
	log.Println("postgres db connected SUCCESSFUL")
	defer dbpool.Close()

	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*")
	r.GET("/", site.Home)
	r.GET("/rest", rest.RootIndex)
	r.GET("/rest/todo", rest.TodoList)
	r.GET("/rest/task/offset", rest.GetTaskOffset)
	r.GET("/rest/task/status", rest.GetTaskStatusList)
	r.POST("/rest/task", rest.CreateTask)
	r.PUT("/rest/task/:id", rest.EditTask)

	r.Run(serverHost + ":" + serverPort)

	log.Println("Todo list App started SUCCESSFULL")
}
