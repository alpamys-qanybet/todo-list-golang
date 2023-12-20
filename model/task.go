package model

import (
	"context"
	"database/sql"
	"fmt"
	"todo/db"
)

type Task struct {
	Id          uint16 `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

const StatusCreated string = "created"
const StatusInProgress string = "in_progress"
const StatusPaused string = "paused"
const StatusDone string = "done"
const StatusDeleted string = "deleted"

type Status struct {
	Name string `json:"name"`
}

func GetTaskTotalElements(status string) (uint32, error) {
	conn, err := db.ConnectionPool()
	if err != nil {
		return 0, err
	}

	var result uint32

	sqlQuery := `
		SELECT COUNT(*) AS _c
		FROM task
	`

	if len(status) > 0 {
		sqlQuery += "WHERE status = '" + status + "'"
	} else {
		sqlQuery += "WHERE status <> '" + StatusDeleted + "'" // by default show all and ignore deleted
	}

	fmt.Println(sqlQuery)

	err = conn.QueryRow(context.Background(), sqlQuery).Scan(&result)

	return result, err
}

func GetTaskListByOffset(offset uint16, limit uint8, status string) ([]*Task, error) {
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
		sqlQuery += "WHERE status = '" + status + "'"
	} else {
		sqlQuery += "WHERE status <> '" + StatusDeleted + "'" // by default show all and ignore deleted
	}
	fmt.Println(sqlQuery)

	rows, err := conn.Query(context.Background(), sqlQuery)
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

	return id, nil
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

	return err
}
