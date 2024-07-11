package views

import (
	"database/sql"

	"github.com/thanakize/todoApiDB/sharedinterface"
)

func GetTodos(db *sql.DB) ([]sharedinterface.Todo, error){
	var todos []sharedinterface.Todo
	rows, err := db.Query("SELECT id, title, status FROM todos")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var todo sharedinterface.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func GetTodoById(db *sql.DB, id string) (sharedinterface.Todo, error){
	var todo sharedinterface.Todo
	row := db.QueryRow("SELECT id, title, status FROM todos where id=$1",id)
	err := row.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		return sharedinterface.Todo{}, err
	}
	return todo, nil
}
func InsertTodo(db *sql.DB, todo sharedinterface.Todo) (sharedinterface.Todo, error){
	q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id, title, status"
	row := db.QueryRow(q, todo.Title, todo.Status)
	var newTodo sharedinterface.Todo
	err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)
	if err != nil {
	
		return sharedinterface.Todo{}, err
	}
	return newTodo, nil
}

func UpdateTodo(db *sql.DB, todo sharedinterface.Todo, id string) (sharedinterface.Todo, error){
	q := "update todos set title = $1, status = $2 where id = $3 RETURNING id, title, status"
	// q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
	row := db.QueryRow(q, todo.Title, todo.Status, id)
	var newTodo sharedinterface.Todo 
	err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

	if err != nil {

		return sharedinterface.Todo{}, err
	}
	return newTodo, nil
}
func DeleteTodo(db *sql.DB, id string) error{
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			return err
		}
	return nil
}
func PatchStatusTodo(db *sql.DB, todo sharedinterface.Todo, id string) (sharedinterface.Todo, error){
	q := "update todos set status = $1 where id = $2 RETURNING id, title, status"
		row := db.QueryRow(q, todo.Status, id) 
		var newTodo sharedinterface.Todo 
		err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

		if err != nil {

			return sharedinterface.Todo{}, err
		}

		return newTodo, nil
}
func PatchTitleTodo(db *sql.DB, todo sharedinterface.Todo, id string) (sharedinterface.Todo, error){
	q := "update todos set title = $1 where id = $2 RETURNING id, title, status"
		// q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
		row := db.QueryRow(q, todo.Title, id)
		var newTodo sharedinterface.Todo 
		err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

		if err != nil {

			return sharedinterface.Todo{}, err
		}

		return newTodo, nil
}
