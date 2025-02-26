package middleware

import (
    "encoding/json"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/gorilla/mux"
)

type contextKey string

const UserIDKey contextKey = "userID"

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
    createdTodo, err := models.CreateTodo(todo.Task, todo.Description, todo.Important, userID)
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
    updatedTodo, err := models.UpdateTodo(params["id"], todo.Task, todo.Description, todo.Done, todo.Important, userID)
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
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    user, err := models.GetUserByUsername(request.Username)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    err = models.ShareTodoWithUser(request.TaskID, user.ID, userID)
    if err != nil {
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
    sharedByMeTodos, err := models.GetSharedByMeTodos(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := map[string]interface{}{
        "received": sharedTodos,
        "shared":   sharedByMeTodos,
    }
    json.NewEncoder(w).Encode(response)
}

// Team-related handlers

// Create a new team
func CreateTeam(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var team models.Team
    _ = json.NewDecoder(r.Body).Decode(&team)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    createdTeam, err := models.CreateTeam(team.Name, team.Password, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTeam)
}

// Join an existing team
func JoinTeam(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        TeamName string `json:"teamName"`
        Password string `json:"password"`
    }
    _ = json.NewDecoder(r.Body).Decode(&request)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    err := models.JoinTeam(request.TeamName, request.Password, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Create a new team todo
func CreateTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var todo models.TeamTodo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    params := mux.Vars(r)
    createdTodo, err := models.CreateTeamTodo(todo.Task, todo.Description, todo.Important, params["teamId"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTodo)
}

// Get all teams for the authenticated user
func GetTeams(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    teams, err := models.GetTeams(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(teams)
}

// Get all team todos
func GetTeamTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    todos, err := models.GetTeamTodos(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(todos)
}

// Update a team todo
func UpdateTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    var todo models.TeamTodo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    updatedTodo, err := models.UpdateTeamTodo(params["id"], todo.Task, todo.Description, todo.Done, todo.Important, todo.TeamID, todo.AssignedTo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(updatedTodo)
}

// Delete a team todo
func DeleteTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    err := models.DeleteTeamTodo(params["id"], params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Remove a team member
func RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    err := models.RemoveTeamMember(params["teamId"], params["userId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Get team details
func GetTeamDetails(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    team, err := models.GetTeamByID(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    todos, err := models.GetTeamTodos(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := map[string]interface{}{
        "team":  team,
        "tasks": todos,
    }
    json.NewEncoder(w).Encode(response)
}

// Get team members
func GetTeamMembers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    members, err := models.GetTeamMembers(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(members)
}

// Add a team member
func AddTeamMember(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        Username string `json:"username"`
    }
    _ = json.NewDecoder(r.Body).Decode(&request)
    params := mux.Vars(r)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    err := models.AddTeamMember(params["teamId"], request.Username, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}


