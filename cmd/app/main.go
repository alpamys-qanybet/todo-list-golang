package main

import (
	"log" // swagger embed files	"log"
	"os"
	"strconv"
	"todo/api"
	"todo/internal/config"
	"todo/pkg/db"

	_ "todo/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var databaseUrl string

func readEnvVariables() (serverHost, serverPort string) {
	_ = godotenv.Load()
	// ignore .env errors, use default values instead
	// in docker specify env variables by ENV command

	serverHost = os.Getenv("SERVER_HOST")

	serverPort = os.Getenv("SERVER_PORT")
	if "" == serverPort {
		serverPort = "8080" // default value
	}

	databaseUrl = os.Getenv("DATABASE_URL")
	if "" == databaseUrl {
		databaseUrl = "postgresql://postgres:postgres@localhost:5432/todo" // default value
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if "" == jwtSecret {
		// log.Fatal("JWT_SECRET is not set in .env file")
	}
	config.SetJwtSecret(jwtSecret)

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	config.SetDebugLog(debug)

	if debug {
		log.Println("environtment variables are read")
	}
	return
}

func connectDB() (*pgxpool.Pool, error) {
	dbpool, err := db.Connect(databaseUrl)
	if err != nil {
		return nil, err
	}

	if config.DebugLog() {
		log.Println("postgres db connected SUCCESSFUL")
	}

	return dbpool, nil
}

func setupRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Recovery()) // recovery middleware

	r.GET("/api", api.RootIndex)
	r.POST("/api/user/login", api.UserLogin)
	r.GET("/api/task", api.GetTaskList)
	r.GET("/api/task/status", api.GetTaskStatusList)
	r.POST("/api/task", api.CreateTask)
	r.GET("/api/task/:id", api.GetTask)
	r.PUT("/api/task/:id", api.EditTask)
	r.PUT("/api/task/:id/start_progress", api.StartTaskProgress)
	r.PUT("/api/task/:id/pause", api.PauseTask)
	r.PUT("/api/task/:id/done", api.DoneTask)
	r.DELETE("/api/task/:id", api.DeleteTask) // only changes status to 'deleted'
	r.PUT("/api/task/:id/restore", api.RestoreTask)
	r.DELETE("/api/task/:id/completely", api.DeleteTaskCompletely)
	r.DELETE("/api/task/free_trash", api.FreeTaskTrash)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return
}

// @title           Todo app API
// @version         1.0
// @description     This is a sample todo app.

// @contact.name   Alpamys Kanibetov
// @contact.email  alpamys.kanibetov@gmail.com

// @host      localhost:8080
// @BasePath  /rest
// @query.collection.format multi

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @query.collection.format multi

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(gin.ReleaseMode)

	serverHost, serverPort := readEnvVariables()

	if config.DebugLog() {
		log.Println("Todo list App init")
	}

	dbpool, err := connectDB()
	if err != nil {
		if config.DebugLog() {
			log.Fatalf("Error on postgres database: %v\n", err)
		}
	}
	defer dbpool.Close()

	r := setupRouter()

	if config.DebugLog() {
		log.Println("Todo list App started SUCCESSFULL")
	}

	r.Run(serverHost + ":" + serverPort)
}
