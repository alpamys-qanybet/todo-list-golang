package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"todo/config"
	"todo/model"
	"todo/rest"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var dbpool *pgxpool.Pool
var r *gin.Engine
var id uint16 // created task id to further tests
var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfY29udGVudCI6MSwiX3RpbWUiOjE3MDMxNzU3NDQ3MTIsIl90b2tlbl9pZCI6ImExMTY3MmE1LTdhNmEtNGZiMi05MDIwLWVkMGEwMTNmZjI4OCJ9.iqzPRTxJwXv6OXlCE4RslhEIPvDUJbwpqSWpB2mY2Uw"
var err error

func TestNewTask(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	_, _ = readEnvVariables()
	config.SetDebugLog(false)

	dbpool, err = ConnectDB()
	if err != nil {
		if config.DebugLog() {
			log.Fatalf("Error on postgres database: %v\n", err)
		}
	}

	r = SetupRouter()
	bodyData := map[string]interface{}{
		"name": "New Task",
	}

	jsonValue, _ := json.Marshal(bodyData)
	req, _ := http.NewRequest("POST", "/rest/task?"+rest.AppSecretName+"="+rest.AppSecret(), bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	id = uint16(result["id"].(float64))
}

func TestGetTask(t *testing.T) {
	req, _ := http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, id, uint16(result["id"].(float64)))
	assert.Equal(t, "New Task", result["name"].(string))
	assert.Equal(t, model.StatusCreated, result["status"].(string))
}

func TestEditTask(t *testing.T) {
	bodyData := map[string]interface{}{
		"name":        "Changed Task",
		"description": "Lorem ipsum",
	}

	jsonValue, _ := json.Marshal(bodyData)
	req, _ := http.NewRequest("PUT",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		bytes.NewBuffer(jsonValue),
	)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, "Changed Task", result["name"].(string))
	assert.Equal(t, "Lorem ipsum", result["description"].(string))
}

func TestStartTaskProgress(t *testing.T) {
	req, _ := http.NewRequest("PUT",
		"/rest/task/"+strconv.Itoa(int(id))+"/start_progress?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, model.StatusInProgress, result["status"].(string))
}

func TestPauseTask(t *testing.T) {
	req, _ := http.NewRequest("PUT",
		"/rest/task/"+strconv.Itoa(int(id))+"/pause?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, model.StatusPaused, result["status"].(string))
}
func TestDoneTask(t *testing.T) {
	req, _ := http.NewRequest("PUT",
		"/rest/task/"+strconv.Itoa(int(id))+"/done?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, model.StatusDone, result["status"].(string))
}

func TestDeleteTask(t *testing.T) {
	defer dbpool.Close()

	req, _ := http.NewRequest("DELETE",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)
	req.Header.Add("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		"/rest/task/"+strconv.Itoa(int(id))+"?"+rest.AppSecretName+"="+rest.AppSecret(),
		nil,
	)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(w.Body.String()), &result)

	assert.Equal(t, model.StatusDeleted, result["status"].(string))
}
