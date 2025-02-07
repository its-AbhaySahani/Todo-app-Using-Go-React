package middleware

import (
    "encoding/json"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

var todos []models.Todo

// Get all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// Create a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var todo models.Todo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    todo.ID = uuid.New().String()
    todos = append(todos, todo)
    json.NewEncoder(w).Encode(todo)
}

// Update an existing todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range todos {
        if item.ID == params["id"] {
            todos = append(todos[:index], todos[index+1:]...)
            var todo models.Todo
            _ = json.NewDecoder(r.Body).Decode(&todo)
            todo.ID = params["id"]
            todos = append(todos, todo)
            json.NewEncoder(w).Encode(todo)
            return
        }
    }
    json.NewEncoder(w).Encode(todos)
}

// Delete a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range todos {
        if item.ID == params["id"] {
            todos = append(todos[:index], todos[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(todos)
}

// Undo a todo
func UndoTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range todos {
        if item.ID == params["id"] {
            todos[index].Done = false
            json.NewEncoder(w).Encode(todos[index])
            return
        }
    }
    json.NewEncoder(w).Encode(todos)
}