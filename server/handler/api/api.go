package api

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/dgrijalva/jwt-go"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/users"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/teams"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/team_members"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/team_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/shared_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/middleware"
)

var jwtKey = []byte("ZLR+ZInOHXQst1seVlV6JVuZe1k3vasV1BRyqAHAyaY=")

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    Username string `json:"username"`
    UserID   string `json:"user_id"`
    jwt.StandardClaims
}

// Register handles user registration
func Register(userService *users.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        res, err := userService.CreateUser(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

// Login handles user authentication and issues JWT tokens
func Login(userService *users.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var creds Credentials
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        user, err := userService.GetUserByUsername(context.Background(), creds.Username)
        if err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        if err := userService.VerifyPassword(user.Password, creds.Password); err != nil {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }

        // Set expiration time for token
        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Username: creds.Username,
            UserID:   user.ID,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        // Create the token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Error generating token", http.StatusInternalServerError)
            return
        }

        // Send the token in response
        json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
    }
}

// Todo Handlers

func GetTodos(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        userID := r.Context().Value(middleware.UserIDKey).(string)
        res, err := todoService.GetTodosByUserID(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func CreateTodo(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the user ID from the context
        userID := r.Context().Value(middleware.UserIDKey).(string)
        req.UserID = userID
        
        res, err := todoService.CreateTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func UpdateTodo(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.UpdateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the user ID from the context and ID from the URL params
        userID := r.Context().Value(middleware.UserIDKey).(string)
        req.UserID = userID
        req.ID = mux.Vars(r)["id"]
        
        res, err := todoService.UpdateTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func DeleteTodo(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        res, err := todoService.DeleteTodo(context.Background(), params["id"], userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func UndoTodo(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        res, err := todoService.UndoTodo(context.Background(), params["id"], userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

// Shared Todos Handlers

func CreateSharedTodo(sharedTodoService *shared_todos.SharedTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateSharedTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the sharedBy field from the context
        sharedBy := r.Context().Value(middleware.UserIDKey).(string)
        req.SharedBy = sharedBy
        
        res, err := sharedTodoService.CreateSharedTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func GetSharedTodos(sharedTodoService *shared_todos.SharedTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        // Get todos shared with the user
        received, err := sharedTodoService.GetSharedTodos(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Get todos shared by the user
        shared, err := sharedTodoService.GetSharedByMeTodos(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Combine the responses
        response := dto.SharedTodosResponse{
            Received: received.Received,
            Shared: shared.Shared,
        }
        
        json.NewEncoder(w).Encode(response)
    }
}

// Team Handlers

func CreateTeam(teamService *teams.TeamService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateTeamRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the admin ID from the context
        adminID := r.Context().Value(middleware.UserIDKey).(string)
        req.AdminID = adminID
        
        res, err := teamService.CreateTeam(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func GetTeams(teamService *teams.TeamService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        res, err := teamService.GetTeamsByAdminID(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func GetTeamByID(teamService *teams.TeamService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        
        res, err := teamService.GetTeamByID(context.Background(), params["teamId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

// Team Todos Handlers

func GetTeamTodos(teamTodoService *team_todos.TeamTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        
        res, err := teamTodoService.GetTeamTodos(context.Background(), params["teamId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func CreateTeamTodo(teamTodoService *team_todos.TeamTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateTeamTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the team ID from the URL parameters
        params := mux.Vars(r)
        req.TeamID = params["teamId"]
        
        res, err := teamTodoService.CreateTeamTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func UpdateTeamTodo(teamTodoService *team_todos.TeamTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.UpdateTeamTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the ID and team ID from the URL parameters
        params := mux.Vars(r)
        req.ID = params["id"]
        req.TeamID = params["teamId"]
        
        res, err := teamTodoService.UpdateTeamTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func DeleteTeamTodo(teamTodoService *team_todos.TeamTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        
        res, err := teamTodoService.DeleteTeamTodo(context.Background(), params["id"], params["teamId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

// Team Members Handlers

func GetTeamMembers(teamMemberService *team_members.TeamMemberService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        
        res, err := teamMemberService.GetTeamMembers(context.Background(), params["teamId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func AddTeamMember(teamMemberService *team_members.TeamMemberService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.AddTeamMemberRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the team ID from the URL parameters
        params := mux.Vars(r)
        req.TeamID = params["teamId"]
        
        res, err := teamMemberService.AddTeamMember(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}

func RemoveTeamMember(teamMemberService *team_members.TeamMemberService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        
        res, err := teamMemberService.RemoveTeamMember(context.Background(), params["teamId"], params["userId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
}