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
    router.HandleFunc("/api/todos", middleware.GetTodos).Methods("GET")
    router.HandleFunc("/api/todo", middleware.CreateTodo).Methods("POST")
    router.HandleFunc("/api/todo/{id}", middleware.UpdateTodo).Methods("PUT")
    router.HandleFunc("/api/todo/{id}", middleware.DeleteTodo).Methods("DELETE")
    router.HandleFunc("/api/todo/undo/{id}", middleware.UndoTodo).Methods("PUT")
    return router
}