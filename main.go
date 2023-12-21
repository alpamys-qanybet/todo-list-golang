package main

import (
	"log"
	"os"
	"strconv"
	"todo/config"
	"todo/db"
	"todo/rest"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
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

	appSecret := os.Getenv("APP_SECRET")
	if "" == appSecret {
		log.Fatal("APP_SECRET is not set in .env file")
	}
	rest.SetAppSecret(appSecret)

	jwtSecret := os.Getenv("JWT_SECRET")
	if "" == jwtSecret {
		log.Fatal("JWT_SECRET is not set in .env file")
	}
	config.SetJwtSecret(jwtSecret)

	debugStr := os.Getenv("DEBUG")
	debug := true
	if "" != debugStr {
		boolValue, err := strconv.ParseBool(debugStr)
		if err == nil {
			debug = boolValue
		} // error occurred, nevermind
	}
	config.SetDebugLog(debug)

	log.Println("environtment variables are read")
	return
}

func ConnectDB() (*pgxpool.Pool, error) {
	dbpool, err := db.Connect(databaseUrl)
	if err != nil {
		return nil, err
	}

	err = db.CreateDatabaseTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	if config.DebugLog() {
		log.Println("postgres db connected SUCCESSFUL")
	}

	return dbpool, nil
}

func SetupRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Recovery()) // recovery middleware

	r.GET("/rest", rest.RootIndex)
	r.POST("/rest/user/login", rest.UserLogin)
	r.GET("/rest/task/offset", rest.GetTaskOffset)
	r.GET("/rest/task/status", rest.GetTaskStatusList)
	r.POST("/rest/task", rest.CreateTask)
	r.GET("/rest/task/:id", rest.GetTask)
	r.PUT("/rest/task/:id", rest.EditTask)
	r.PUT("/rest/task/:id/start_progress", rest.StartTaskProgress)
	r.PUT("/rest/task/:id/pause", rest.PauseTask)
	r.PUT("/rest/task/:id/done", rest.DoneTask)
	r.DELETE("/rest/task/:id", rest.DeleteTask) // only changes status to 'deleted'
	r.PUT("/rest/task/:id/restore", rest.RestoreTask)
	r.DELETE("/rest/task/:id/completely", rest.DeleteTaskCompletely)
	r.DELETE("/rest/task/free_trash", rest.FreeTaskTrash)

	return
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	serverHost, serverPort := readEnvVariables()

	if config.DebugLog() {
		log.Println("Todo list App init")
	}

	dbpool, err := ConnectDB()
	if err != nil {
		if config.DebugLog() {
			log.Fatalf("Error on postgres database: %v\n", err)
		}
	}
	defer dbpool.Close()

	r := SetupRouter()

	if config.DebugLog() {
		log.Println("Todo list App started SUCCESSFULL")
	}

	r.Run(serverHost + ":" + serverPort)
}
