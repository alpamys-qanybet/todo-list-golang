package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"todo/controller"

	"github.com/gin-gonic/gin"
)

func GetTaskOffset(c *gin.Context) {
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

	fmt.Println(offset, limit)

	res, err := controller.GetTaskOffset(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("res")
	fmt.Println(res)

	c.JSON(http.StatusOK, res)
}

func GetTaskStatusList(c *gin.Context) {
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
	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid body param(s)")
		return
	}

	name := bodyData["name"].(string)

	id, err := controller.CreateTask(name)
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
