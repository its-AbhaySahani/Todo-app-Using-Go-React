package handler

import (
    "database/sql"
    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/api"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/users_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/teams_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_members_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/shared_todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/users"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/teams"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/team_members"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/team_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Services/shared_todos"
)

func SetupRoutes(router *mux.Router, DB *sql.DB) {
    // Initialize repositories
    userRepo := users_repository.NewUserRepository(DB)
    todoRepo := todos_repository.NewTodoRepository(DB)
    teamRepo := teams_repository.NewTeamRepository(DB)
    teamMemberRepo := team_members_repository.NewTeamMemberRepository(DB)
    teamTodoRepo := team_todos_repository.NewTeamTodoRepository(DB)
    sharedTodoRepo := shared_todos_repository.NewSharedTodoRepository(DB)

    // Initialize services
    userService := users.NewUserService(userRepo)
    todoService := todos.NewTodoService(todoRepo)
    teamService := teams.NewTeamService(teamRepo)
    teamMemberService := team_members.NewTeamMemberService(teamMemberRepo)
    teamTodoService := team_todos.NewTeamTodoService(teamTodoRepo)
    sharedTodoService := shared_todos.NewSharedTodoService(sharedTodoRepo)

    // User routes
    router.HandleFunc("/api/register", api.Register(userService)).Methods("POST")
    router.HandleFunc("/api/login", api.Login(userService)).Methods("POST")

    // Todo routes
    router.HandleFunc("/api/todos", api.GetTodos(todoService)).Methods("GET")
    router.HandleFunc("/api/todo", api.CreateTodo(todoService)).Methods("POST")
    router.HandleFunc("/api/todo/{id}", api.UpdateTodo(todoService)).Methods("PUT")
    router.HandleFunc("/api/todo/{id}", api.DeleteTodo(todoService)).Methods("DELETE")
    router.HandleFunc("/api/todo/undo/{id}", api.UndoTodo(todoService)).Methods("PUT")
    // router.HandleFunc("/api/share", api.ShareTodo(todoService)).Methods("POST")
    router.HandleFunc("/api/shared", api.GetSharedTodos(sharedTodoService)).Methods("GET")

    // Team routes
    router.HandleFunc("/api/team", api.CreateTeam(teamService)).Methods("POST")
    // router.HandleFunc("/api/team/join", api.JoinTeam(teamService)).Methods("POST")
    router.HandleFunc("/api/teams", api.GetTeams(teamService)).Methods("GET")
    // router.HandleFunc("/api/team/{teamId}", api.GetTeamDetails(teamService)).Methods("GET")
    router.HandleFunc("/api/team/{teamId}/todos", api.GetTeamTodos(teamTodoService)).Methods("GET")
    router.HandleFunc("/api/team/{teamId}/todo", api.CreateTeamTodo(teamTodoService)).Methods("POST")
    router.HandleFunc("/api/team/{teamId}/todo/{id}", api.UpdateTeamTodo(teamTodoService)).Methods("PUT")
    router.HandleFunc("/api/team/{teamId}/todo/{id}", api.DeleteTeamTodo(teamTodoService)).Methods("DELETE")
    router.HandleFunc("/api/team/{teamId}/members", api.GetTeamMembers(teamMemberService)).Methods("GET")
    router.HandleFunc("/api/team/{teamId}/member", api.AddTeamMember(teamMemberService)).Methods("POST")
    router.HandleFunc("/api/team/{teamId}/member/{userId}", api.RemoveTeamMember(teamMemberService)).Methods("DELETE")
}