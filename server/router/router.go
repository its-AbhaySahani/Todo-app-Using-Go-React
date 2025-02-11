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

    return router
}