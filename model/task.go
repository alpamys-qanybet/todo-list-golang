package model

import (
	"context"
	"fmt"
	"todo/db"
)

type Task struct {
	Id   uint16 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Status			string	`json:"status,omitempty"`
	// Description		string	`json:"description,omitempty"`
}

const StatusCreated string = "created"
const StatusInProgress string = "in_progress"
const StatusPaused string = "paused"
const StatusDone string = "done"
const StatusDeleted string = "deleted"

type Status struct {
	Name string `json:"name,omitempty"`
}

func GetTaskTotalElements() (uint32, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return 0, err
	}

	var result uint32

	err = conn.QueryRow(context.Background(), `
        SELECT COUNT(*) AS _c
        FROM task
    `).Scan(&result)

	return result, err
}

func GetTaskListByOffset(offset uint16, limit uint8) ([]*Task, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	var result []*Task = []*Task{}

	rows, err := conn.Query(context.Background(), `
        SELECT id, name
        FROM task
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// var id uint16
		// var name string
		var item Task

		// err = rows.Scan(&id, &name)
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}
	err = rows.Err()

	fmt.Println("returning result")
	fmt.Println(result)
	return result, err
}

func GetStatusList() ([]*Status, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	var result []*Status = []*Status{}

	rows, err := conn.Query(context.Background(), `
        SELECT name
        FROM status
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Status

		err = rows.Scan(&item.Name)
		if err != nil {
			return nil, err
		}

		result = append(result, &item)
	}
	err = rows.Err()

	return result, err
}

func CreateTask(name string) (uint16, error) {
	var id uint16

	conn, err := db.ConnectionPool()
	if err != nil {
		return id, err
	}

	err = conn.QueryRow(context.Background(), `
        INSERT INTO task(name)
        VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}
