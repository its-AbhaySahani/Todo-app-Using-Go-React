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
    userID := r.Context().Value("userID").(string)
    todos, err := models.GetTodos(userID)
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
    userID := r.Context().Value("userID").(string)
    createdTodo, err := models.CreateTodo(todo.Task, userID)
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
    userID := r.Context().Value("userID").(string)
    updatedTodo, err := models.UpdateTodo(params["id"], todo.Task, todo.Done, userID)
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
    userID := r.Context().Value("userID").(string)
    err := models.DeleteTodo(params["id"], userID)
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
    userID := r.Context().Value("userID").(string)
    undoneTodo, err := models.UndoTodo(params["id"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(undoneTodo)
}