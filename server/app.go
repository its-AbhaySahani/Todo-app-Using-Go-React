package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/infra"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

// RequestLogger is middleware that logs the details of each request
func RequestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
        fmt.Printf("[%s] Completed %s %s in %v\n", 
            time.Now().Format("2006-01-02 15:04:05"), 
            r.Method, 
            r.URL.Path, 
            time.Since(startTime))
    })
}

// PrintRoutes walks through all registered routes and prints them
func PrintRoutes(r *mux.Router) {
    fmt.Println("\n📋 Available API endpoints:")
    fmt.Println("=======================")
    
    err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        pathTemplate, err := route.GetPathTemplate()
        if err != nil {
            return nil // Skip routes without a path template
        }
        
        methods, err := route.GetMethods()
        if err != nil {
            methods = []string{"ANY"} // Default if methods are not specified
        }
        
        fmt.Printf("%-7s %s\n", strings.Join(methods, ","), pathTemplate)
        return nil
    })
    
    if err != nil {
        fmt.Println("Error walking routes:", err)
    }
    fmt.Println("=======================")
}

func main() {
    // Connect to the database
    infra.Connect()

    // Create a new router
    router := mux.NewRouter()

    // Setup routes
    handler.SetupRoutes(router, infra.DB)
    
    // Add request logging middleware to all routes
    router.Use(RequestLogger)
    
    // Print all registered routes
    PrintRoutes(router)

    // Configure CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    // Start the server
    handler := c.Handler(router)
    fmt.Println("\n🚀 Starting the server on port 9000...")
    log.Fatal(http.ListenAndServe(":9000", handler))
}