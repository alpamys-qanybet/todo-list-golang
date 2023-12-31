package model

import (
	"context"
	"database/sql"
	"log"
	"todo/internal/config"
	"todo/pkg/db"
)

type Task struct {
	Id          uint16 `json:"id"`
	Name        string `json:"name" example:"New Task"`
	Status      string `json:"status"`
	Description string `json:"description" example:"Lorum ipsum"`
}

const (
	StatusCreated    = "created"
	StatusInProgress = "in_progress"
	StatusPaused     = "paused"
	StatusDone       = "done"
	StatusDeleted    = "deleted"
)

type Status struct {
	Name string `json:"name"`
}

func GetTaskList(status string) ([]*Task, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	var result []*Task = []*Task{}

	sqlQuery := `
		SELECT id, name, description, status
		FROM task
	`

	if len(status) > 0 {
		sqlQuery += "WHERE status = $1"
	} else {
		sqlQuery += "WHERE status <> $1" // by default show all and ignore deleted
		status = StatusDeleted
	}

	sqlQuery += " ORDER BY id ASC"

	rows, err := conn.Query(context.Background(), sqlQuery, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Task
		var description sql.NullString

		err = rows.Scan(&item.Id, &item.Name, &description, &item.Status)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			if len(description.String) > 0 {
				item.Description = description.String
			}
		}

		result = append(result, &item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if config.DebugLog() {
		log.Println("task list offset: successfully retrived data from db")
	}

	return result, nil
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
	if err != nil {
		return nil, err
	}

	if config.DebugLog() {
		log.Println("task status list: successfully retrived data from db")
	}

	return result, nil

}

func CreateTask(name, description string) (uint16, error) {
	var id uint16

	conn, err := db.ConnectionPool()
	if err != nil {
		return id, err
	}

	err = conn.QueryRow(context.Background(), `
		INSERT INTO task(name, description, status)
		VALUES ($1, $2, $3) RETURNING id`,
		name,
		description,
		StatusCreated,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	if config.DebugLog() {
		log.Println("task create: successfully created data in db")
	}

	return id, nil
}

func GetTask(id uint16) (*Task, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return nil, err
	}

	var item Task
	var description sql.NullString

	err = conn.QueryRow(context.Background(), `
		SELECT id, name, description, status
		FROM task
		WHERE id = $1
	`, id).Scan(&item.Id, &item.Name, &description, &item.Status)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		if len(description.String) > 0 {
			item.Description = description.String
		}
	}

	if config.DebugLog() {
		log.Println("get task: successfully retrieved data in db")
	}

	return &item, nil
}

func EditTask(id uint16, name, description string) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET name = $1,
			description = $2
		WHERE id = $3`,
		name,
		description,
		id,
	)
	if err != nil {
		return err
	}

	if config.DebugLog() {
		log.Println("task edit: successfully edited data in db")
	}

	return nil
}

func StartTaskProgress(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET status = $1
		WHERE id = $2`,
		StatusInProgress,
		id,
	)

	if config.DebugLog() {
		log.Println("task start progress: successfully changed status in db")
	}

	return err
}

func PauseTask(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET status = $1
		WHERE id = $2`,
		StatusPaused,
		id,
	)

	if config.DebugLog() {
		log.Println("task pause: successfully changed status in db")
	}

	return err
}

func DoneTask(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET status = $1
		WHERE id = $2`,
		StatusDone,
		id,
	)

	if config.DebugLog() {
		log.Println("task done: successfully changed status in db")
	}

	return err
}

func DeleteTask(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET status = $1
		WHERE id = $2`,
		StatusDeleted,
		id,
	)

	if config.DebugLog() {
		log.Println("delete task: successfully deleted task in db")
	}

	return err
}

func RestoreTask(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		UPDATE task
		SET status = $1
		WHERE id = $2`,
		StatusCreated, // ? maybe we should use in progress status or paused
		id,
	)

	if config.DebugLog() {
		log.Println("restore task: successfully restored task in db")
	}

	return err
}

func DeleteTaskCompletely(id uint16) error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		DELETE FROM task
		WHERE id = $1`,
		id,
	)

	if config.DebugLog() {
		log.Println("delete task completely: successfully deleted task in db completely")
	}

	return err
}

func FreeTaskTrash() error {
	conn, err := db.ConnectionPool()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		DELETE FROM task
		WHERE status = $1`,
		StatusDeleted,
	)

	if config.DebugLog() {
		log.Println("free task trash: successfully freed task trash in db")
	}

	return err
}
