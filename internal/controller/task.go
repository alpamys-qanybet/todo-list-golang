package controller

import (
	"errors"
	"log"
	"strings"
	"todo/internal/config"
	"todo/internal/model"
)

func GetTaskList(status string) (interface{}, error) {
	if len(status) > 0 {
		if !(status == model.StatusCreated || status == model.StatusInProgress || status == model.StatusPaused || status == model.StatusDone || status == model.StatusDeleted) {
			// it is not one of our statuses, user just mistyped something else
			if config.DebugLog() {
				log.Printf("task offset incorrect status: %s", status)
			}
			status = ""
		}
	}

	list, err := model.GetTaskList(status)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetTaskStatusList() (interface{}, error) {
	list, err := model.GetStatusList()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func CreateTask(name, description string) (uint16, error) {
	if strings.Trim(name, " ") == "" {
		return uint16(0), errors.New("create_task_failure_name_is_required")
	}

	return model.CreateTask(name, description)
}

func GetTask(id uint16) (*model.Task, error) {
	return model.GetTask(id)
}

func EditTask(id uint16, name, description string) error {
	if strings.Trim(name, " ") == "" {
		return errors.New("edit_task_failure_name_is_required")
	}

	return model.EditTask(id, name, description)
}

func StartTaskProgress(id uint16) error {
	return model.StartTaskProgress(id)
}

func PauseTask(id uint16) error {
	return model.PauseTask(id)
}

func DoneTask(id uint16) error {
	return model.DoneTask(id)
}

func DeleteTask(id uint16) error {
	return model.DeleteTask(id)
}

func RestoreTask(id uint16) error {
	return model.RestoreTask(id)
}

func DeleteTaskCompletely(id uint16) error {
	return model.DeleteTaskCompletely(id)
}

func FreeTaskTrash() error {
	return model.FreeTaskTrash()
}
