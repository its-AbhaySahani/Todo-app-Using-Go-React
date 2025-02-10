package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/router"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/rs/cors"
)

func main() {
    database.Connect()
    r := router.Router()
    fmt.Println("starting the server on port 9000...")

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)
    log.Fatal(http.ListenAndServe(":9000", handler))
}