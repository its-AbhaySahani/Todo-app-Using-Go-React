package middleware

import (
    "encoding/json"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models"
    "github.com/gorilla/mux"
)

// Get all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    todos, err := models.GetTodos()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(todos)
}

// Create a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var todo models.Todo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    createdTodo, err := models.CreateTodo(todo.Task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTodo)
}

// Update an existing todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    var todo models.Todo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    updatedTodo, err := models.UpdateTodo(params["id"], todo.Task, todo.Done)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(updatedTodo)
}

// Delete a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    err := models.DeleteTodo(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Undo a todo
func UndoTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    undoneTodo, err := models.UndoTodo(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(undoneTodo)
}