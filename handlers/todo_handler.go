package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"kelabang/models"
)

func GetTodos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.URL)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, task FROM todos ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.Task == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	insertTodo := "INSERT INTO todos (task) VALUES ($1) RETURNING id"
	err = db.QueryRow(insertTodo, todo.Task).Scan(&todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/todos/"):])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.Task == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	updateTodo := "UPDATE todos SET task = $1 WHERE id = $2"
	_, err = db.Exec(updateTodo, todo.Task, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/todos/"):])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	deleteTodo := "DELETE FROM todos WHERE id = $1"
	_, err = db.Exec(deleteTodo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
