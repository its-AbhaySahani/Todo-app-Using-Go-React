package router

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
)

func Router() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, server!")
    }).Methods("GET")

    router.HandleFunc("/api/register", middleware.Register).Methods("POST")
    router.HandleFunc("/api/login", middleware.Login).Methods("POST")

    apiRouter := router.PathPrefix("/api").Subrouter()
    apiRouter.Use(middleware.AuthMiddleware)
    apiRouter.HandleFunc("/todos", middleware.GetTodos).Methods("GET")
    apiRouter.HandleFunc("/todo", middleware.CreateTodo).Methods("POST")
    apiRouter.HandleFunc("/todo/{id}", middleware.UpdateTodo).Methods("PUT")
    apiRouter.HandleFunc("/todo/{id}", middleware.DeleteTodo).Methods("DELETE")
    apiRouter.HandleFunc("/todo/undo/{id}", middleware.UndoTodo).Methods("PUT")
    apiRouter.HandleFunc("/share", middleware.ShareTodo).Methods("POST")
    apiRouter.HandleFunc("/shared", middleware.GetSharedTodos).Methods("GET")

    // Team routes
    apiRouter.HandleFunc("/team", middleware.CreateTeam).Methods("POST")
    apiRouter.HandleFunc("/team/join", middleware.JoinTeam).Methods("POST")
    apiRouter.HandleFunc("/teams", middleware.GetTeams).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}", middleware.GetTeamDetails).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/todos", middleware.GetTeamTodos).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/todo", middleware.CreateTeamTodo).Methods("POST")
    apiRouter.HandleFunc("/team/{teamId}/todo/{id}", middleware.UpdateTeamTodo).Methods("PUT")
    apiRouter.HandleFunc("/team/{teamId}/todo/{id}", middleware.DeleteTeamTodo).Methods("DELETE")
    apiRouter.HandleFunc("/team/{teamId}/members", middleware.GetTeamMembers).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/member", middleware.AddTeamMember).Methods("POST")
    apiRouter.HandleFunc("/team/{teamId}/member/{userId}", middleware.RemoveTeamMember).Methods("DELETE")

    return router
}