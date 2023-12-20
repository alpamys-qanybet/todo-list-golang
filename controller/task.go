package controller

import (
	"errors"
	"strings"
	"todo/model"
)

func GetTaskOffset(offset uint16, limit uint8) (interface{}, error) {
	totalElements, err := model.GetTaskTotalElements()
	if err != nil {
		return nil, err
	}

	list, err := model.GetTaskListByOffset(offset, limit)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"totalElements": totalElements,
		"list":          list,
	}

	return data, nil
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

func EditTask(id uint16, name, description string) error {
	if strings.Trim(name, " ") == "" {
		return errors.New("edit_task_failure_name_is_required")
	}

	return model.EditTask(id, name, description)
}
