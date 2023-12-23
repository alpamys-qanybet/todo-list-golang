package main

import (
	"log" // swagger embed files	"log"
	"os"
	"strconv"
	"todo/config"
	"todo/db"
	"todo/rest"

	_ "todo/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var databaseUrl string

func readEnvVariables() (serverHost, serverPort string) {
	_ = godotenv.Load() // do nothing if file .env does not exist, just use default values.

	serverHost = os.Getenv("SERVER_HOST")

	serverPort = os.Getenv("SERVER_PORT")
	if "" == serverPort {
		serverPort = "8080"
	}

	databaseUrl = os.Getenv("DATABASE_URL")
	if "" == databaseUrl {
		databaseUrl = "postgresql://postgres:postgres@localhost:5432/todo"
	}

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

	if debug {
		log.Println("environtment variables are read")
	}
	return
}

func ConnectDB() (*pgxpool.Pool, error) {
	dbpool, err := db.Connect(databaseUrl)
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
	r.GET("/rest/task", rest.GetTaskList)
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
