package api

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/dgrijalva/jwt-go"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/users"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/teams"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/team_members"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/team_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/shared_todos"
    persistentDto "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
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

func Register(userService *users.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req persistentDto.CreateUserRequest
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

func Login(userService *users.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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

        err = userService.VerifyPassword(user.Password, creds.Password)
        if err != nil {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }

        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Username: creds.Username,
            UserID:   user.ID,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Error generating token", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })

        json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
    }
}

func GetTodos(todoService *todos.TodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(string)
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
        var req persistentDto.CreateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
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
        var req persistentDto.UpdateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
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
        params := mux.Vars(r)
        userID := r.Context().Value("userID").(string)
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
        params := mux.Vars(r)
        userID := r.Context().Value("userID").(string)
        res, err := todoService.UndoTodo(context.Background(), params["id"], userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}

// func ShareTodo(todoService *todos.TodoService) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         var req persistentDto.ShareTodoRequest
//         if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//             http.Error(w, "Invalid request payload", http.StatusBadRequest)
//             return
//         }
//         userID := r.Context().Value("userID").(string)
//         res, err := todoService.ShareTodoWithUser(context.Background(), req.TaskID, req.Username, userID)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         json.NewEncoder(w).Encode(res)
//     }
// }

func GetSharedTodos(sharedTodoService *shared_todos.SharedTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(string)
        res, err := sharedTodoService.GetSharedTodos(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}

func CreateTeam(teamService *teams.TeamService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req persistentDto.CreateTeamRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        res, err := teamService.CreateTeam(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}

// func JoinTeam(teamService *teams.TeamService) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         var req persistentDto.JoinTeamRequest
//         if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//             http.Error(w, "Invalid request payload", http.StatusBadRequest)
//             return
//         }
//         userID := r.Context().Value("userID").(string)
//         res, err := teamService.JoinTeam(context.Background(), req.TeamName, req.Password, userID)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         json.NewEncoder(w).Encode(res)
//     }
// }

func GetTeams(teamService *teams.TeamService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(string)
        res, err := teamService.GetTeamsByAdminID(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}

// func GetTeamDetails(teamService *teams.TeamService) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         params := mux.Vars(r)
//         res, err := teamService.GetTeamDetails(context.Background(), params["teamId"])
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         json.NewEncoder(w).Encode(res)
//     }
// }

func GetTeamTodos(teamTodoService *team_todos.TeamTodoService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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
        var req persistentDto.CreateTeamTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
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
        var req persistentDto.UpdateTeamTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
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
        params := mux.Vars(r)
        res, err := teamTodoService.DeleteTeamTodo(context.Background(), params["id"], params["teamId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}

func GetTeamMembers(teamMemberService *team_members.TeamMemberService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
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
        var req persistentDto.AddTeamMemberRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
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
        params := mux.Vars(r)
        res, err := teamMemberService.RemoveTeamMember(context.Background(), params["teamId"], params["userId"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(res)
    }
}