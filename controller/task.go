package controller

import (
	"errors"
	"fmt"
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

	fmt.Println("task controller", list)

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

func CreateTask(name string) (uint16, error) {
	if strings.Trim(name, " ") == "" {
		return uint16(0), errors.New("create_task_failure_name_is_required")
	}

	return model.CreateTask(name)
}
