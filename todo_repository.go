package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type TodoRepository struct {
	db *sql.DB
}

func (t TodoRepository) list() ([]Todo, error) {
	var result []Todo

	rows, err := t.db.Query("SELECT * FROM todos")
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var res Todo
		if err := rows.Scan(&res.ID, &res.Description, &res.IsDone); err != nil {
			return []Todo{}, err
		}
		result = append(result, res)
	}
	return result, nil
}

func (t TodoRepository) get(id int64) (Todo, error) {
	var res Todo

	rows, err := t.db.Query("SELECT * FROM todos WHERE id=$1", id)
	if err != nil {
		return res, err
	}
	if !rows.Next() {
		return res, errors.New(fmt.Sprintf("Object with id: %d does not exist.", id))
	}
	if err := rows.Scan(&res.ID, &res.Description, &res.IsDone); err != nil {
		return res, err
	}

	return res, nil
}

func (t TodoRepository) create(description string) (Todo, error) {
	var id int64

	err := t.db.QueryRow("INSERT INTO todos ( description, is_done ) VALUES ( $1, false ) RETURNING id", description).Scan(&id)
	if err != nil {
		return Todo{}, err
	}

	return t.get(id)
}

func (t TodoRepository) delete(id int64) (err error) {
	_, err = t.db.Exec("DELETE from todos WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (t TodoRepository) update(id int64, data Todo) (err error) {
	result, err := t.db.Exec("UPDATE todos SET description=$1, is_done=$2 WHERE id=$3", data.Description, data.IsDone, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		errors.New("Todo not updated.")
	}
	return nil
}
