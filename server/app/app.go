package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/infra"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    // Connect to the database
    infra.Connect()

    // Create a new router
    router := mux.NewRouter()

    // Setup routes
    handler.SetupRoutes(router, infra.DB)

    // Configure CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    // Start the server
    handler := c.Handler(router)
    fmt.Println("Starting the server on port 9000...")
    log.Fatal(http.ListenAndServe(":9000", handler))
}