package middleware

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models"
    "github.com/gorilla/mux"
)

// Get all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
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
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
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
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
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
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
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
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    undoneTodo, err := models.UndoTodo(params["id"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(undoneTodo)
}

// Share a todo with another user
func ShareTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        TaskID   string `json:"taskId"`
        Username string `json:"username"`
    }
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        fmt.Println("Error decoding request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    fmt.Println("Share request:", request)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    user, err := models.GetUserByUsername(request.Username)
    if err != nil {
        fmt.Println("Error fetching user:", err)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    err = models.ShareTodoWithUser(request.TaskID, user.ID, userID)
    if err != nil {
        fmt.Println("Error sharing task:", err)
        http.Error(w, "Error sharing task", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Get shared todos
func GetSharedTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    sharedTodos, err := models.GetSharedTodos(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(sharedTodos)
}