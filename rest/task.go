package rest

import (
	"log"
	"net/http"

	_ "todo/docs"

	"todo/config"
	"todo/controller"

	"github.com/gin-gonic/gin"
)

// GetTaskOffset godoc
// @ID get-task-list
// @Summary      Get task list
// @Description  Get task list
// @Tags         task
// @Accept       json
// @Produce      json
// @Success 200 {array} model.Task
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task [get]
func GetTaskList(c *gin.Context) {
	status := c.Query("status")

	if config.DebugLog() {
		log.Println("requesting task offset", fullUrl(c))
	}

	res, err := controller.GetTaskList(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetTaskStatusList godoc
// @ID get-task-status-list
// @Summary      Get task status list
// @Description  Get task status list
// @Tags         task
// @Accept       json
// @Produce      json
// @Success 200 {array} model.Status
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/status [get]
func GetTaskStatusList(c *gin.Context) {
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

// CreateTask godoc
// @ID create-task
// @Security ApiKeyAuth
// @Summary      Create task
// @Description  Create task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param input body todo.Model true "task input name,description"
// @Success 201 {object} model.Task
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task [post]
func CreateTask(c *gin.Context) {
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

// GetTask godoc
// @ID get-task
// @Summary      Get task
// @Description  Get task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id"
// @Success	200 {object} model.Task
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id} [get]
func GetTask(c *gin.Context) {
	if config.DebugLog() {
		log.Println("requesting task by id", fullUrl(c))
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	res, err := controller.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// EditTask godoc
// @ID edit-task
// @Security ApiKeyAuth
// @Summary      Edit task
// @Description  Edit task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id"
// @Param input body todo.Model true "task input name,description"
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id} [put]
func EditTask(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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

// StartTaskProgress godoc
// @ID start-task-progress
// @Security ApiKeyAuth
// @Summary      Start task progress
// @Description  Start task progress
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id}/start_progress [put]
func StartTaskProgress(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// PauseTask godoc
// @ID pause-task
// @Security ApiKeyAuth
// @Summary      Pause task
// @Description  Pause task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id}/pause [put]
func PauseTask(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// DoneTask godoc
// @ID done-task
// @Security ApiKeyAuth
// @Summary      Done task
// @Description  Done task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id}/done [put]
func DoneTask(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// DeleteTask godoc
// @ID delete-task
// @Security ApiKeyAuth
// @Summary      Delete task
// @Description  Delete task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id} [delete]
func DeleteTask(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// RestoreTask godoc
// @ID restore-task
// @Security ApiKeyAuth
// @Summary      Restore task
// @Description  Restore task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id}/restore [put]
func RestoreTask(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// DeleteTask godoc
// @ID delete-task-completely
// @Security ApiKeyAuth
// @Summary      Delete task completely
// @Description  Delete task completely
// @Tags         task
// @Accept       json
// @Produce      json
// @Param id path int true "task id" default(1)
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/{id}/completely [delete]
func DeleteTaskCompletely(c *gin.Context) {
	_, err := authorizeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := controller.StringToUint16(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

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

// FreeTaskTrash godoc
// @ID free-task-trash
// @Security ApiKeyAuth
// @Summary      Free task trash
// @Description  Free task trash
// @Tags         task
// @Accept       json
// @Produce      json
// @Success	200 {object}
// @Failure      400  {object}  http.StatusBadRequest
// @Failure      500  {object}  http.StatusInternalServerError

// @Router       /task/free_trash [delete]
func FreeTaskTrash(c *gin.Context) {
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
