package rest

import (
	"fmt"
	"net/http"
	"strconv"

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
	fmt.Println("status")
	fmt.Println(status)

	res, err := controller.GetTaskOffset(offset, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("res")
	fmt.Println(res)

	c.JSON(http.StatusOK, res)
}

func GetTaskStatusList(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	res, err := controller.GetTaskStatusList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("res")
	fmt.Println(res)

	c.JSON(http.StatusOK, res)
}

func CreateTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	name := bodyData["name"].(string)
	var description string
	if bodyData["description"] != nil {
		description = bodyData["description"].(string)
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

	c.JSON(http.StatusOK, data)
}

func EditTask(c *gin.Context) {
	if !appSecretIsValid(c) {
		return
	}

	id := controller.StringToUint16(c.Param("id"))
	// fmt.Println(id)

	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	name := bodyData["name"].(string)
	var description string
	if bodyData["description"] != nil {
		description = bodyData["description"].(string)
	}
	// status := bodyData["status"].(string) // implement change status code in another function?

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

	fmt.Println("StartTaskProgress")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.StartTaskProgress(id)
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

	fmt.Println("PauseTask")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.PauseTask(id)
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

	fmt.Println("DoneTask")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.DoneTask(id)
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

	fmt.Println("DeleteTask")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.DeleteTask(id)
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

	fmt.Println("RestoreTask")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.RestoreTask(id)
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

	fmt.Println("DeleteTaskCompletely")
	id := controller.StringToUint16(c.Param("id"))
	fmt.Println(id)

	err := controller.DeleteTaskCompletely(id)
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

	fmt.Println("FreeTaskTrash")

	err := controller.FreeTaskTrash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}
