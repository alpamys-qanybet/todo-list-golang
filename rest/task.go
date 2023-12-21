package rest

import (
	"log"
	"net/http"
	"strconv"

	"todo/config"
	"todo/controller"

	"github.com/gin-gonic/gin"
)

func GetTaskOffset(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	var offset uint16
	var limit uint8
	offsetI, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = uint16(0)
	} else {
		offset = uint16(offsetI)
	}

	limitI, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = uint8(50)
	} else {
		limit = uint8(limitI)
	}

	status := c.Query("status")

	if config.DebugLog() {
		log.Println("requesting task offset", fullUrl(c))
	}

	res, err := controller.GetTaskOffset(offset, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetTaskStatusList(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	if config.DebugLog() {
		log.Println("requesting task status list", fullUrl(c))
	}

	res, err := controller.GetTaskStatusList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	var bodyData map[string]interface{}
	err = extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	name := bodyData["name"].(string)
	var description string
	if bodyData["description"] != nil {
		description = bodyData["description"].(string)
	}

	if config.DebugLog() {
		log.Println("requesting task create", fullUrl(c))
	}

	id, err := controller.CreateTask(name, description)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "create_task_failure_name_is_required" {
			c.JSON(http.StatusUnprocessableEntity, errMsg)
			return
		}
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	data := gin.H{
		"id": id,
	}

	c.JSON(http.StatusCreated, data)
}

func GetTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	if config.DebugLog() {
		log.Println("requesting task by id", fullUrl(c))
	}

	id := controller.StringToUint16(c.Param("id"))

	res, err := controller.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func EditTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	var bodyData map[string]interface{}
	err = extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	name := bodyData["name"].(string)
	var description string
	if bodyData["description"] != nil {
		description = bodyData["description"].(string)
	}

	if config.DebugLog() {
		log.Println("requesting task edit", fullUrl(c))
	}

	err = controller.EditTask(id, name, description)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "edit_task_failure_name_is_required" {
			c.JSON(http.StatusUnprocessableEntity, errMsg)
			return
		}
		c.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func StartTaskProgress(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task start progress", fullUrl(c))
	}

	err = controller.StartTaskProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func PauseTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task pause", fullUrl(c))
	}

	err = controller.PauseTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func DoneTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task done", fullUrl(c))
	}

	err = controller.DoneTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func DeleteTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task delete", fullUrl(c))
	}

	err = controller.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func RestoreTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task restore", fullUrl(c))
	}

	err = controller.RestoreTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func DeleteTaskCompletely(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id := controller.StringToUint16(c.Param("id"))

	if config.DebugLog() {
		log.Println("requesting task delete competely", fullUrl(c))
	}

	err = controller.DeleteTaskCompletely(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}

func FreeTaskTrash(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	if config.DebugLog() {
		log.Println("requesting free task trash", fullUrl(c))
	}

	err = controller.FreeTaskTrash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}
