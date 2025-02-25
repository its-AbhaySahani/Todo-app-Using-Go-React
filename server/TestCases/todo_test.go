package TestCases

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Test User ID that will be used consistently across all tests
const testUserID = "todo-test-user-id"
const testUsername = "todo_testuser" // Different from "testuser" used in auth tests

// Helper function to add a userID to the request context
func addUserIDToContext(r *http.Request, userID string) *http.Request {
    ctx := context.WithValue(r.Context(), "userID", userID)
    return r.WithContext(ctx)
}

// Helper function to ensure the test user exists
func ensureTestUserExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testUserID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if test user exists: %v", err)
    }

    if count == 0 {
        // Create the test user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testUserID, testUsername, "$2a$10$TestHashedPasswordForTestUser",
        )
        if err != nil {
            t.Fatalf("Failed to create test user: %v", err)
        }
        fmt.Println("Created test user with ID:", testUserID)
    } else {
        fmt.Println("Test user already exists with ID:", testUserID)
    }
}

// Helper function to cleanup test data after tests
func cleanupTodoTestData() {
    // Delete todos created by test user
    _, err := database.DB.Exec("DELETE FROM todos WHERE user_id = ?", testUserID)
    if err != nil {
        fmt.Println("Error cleaning up test todos:", err)
    } else {
        fmt.Println("Test todos cleaned up")
    }
    
    // Delete test user
    _, err = database.DB.Exec("DELETE FROM users WHERE id = ?", testUserID)
    if err != nil {
        fmt.Println("Error cleaning up test user:", err)
    } else {
        fmt.Println("Test user cleaned up")
    }
}

// TestGetTodos tests the GetTodos endpoint
func TestGetTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTodos")
    fmt.Println("Testing getting todos for the authenticated user")
    
    // Clean up before test to ensure fresh state
    cleanupTodoTestData()
    
    // Ensure test user exists
    ensureTestUserExists(t)

    // Generate a valid token for testing
    token, err := generateToken(testUsername, testUserID)
    if err != nil {
        t.Fatal("Failed to generate token:", err)
    }

    // Create a new request
    req, err := http.NewRequest("GET", "/api/todos", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    
    // Add authorization header
    req.Header.Set("Authorization", "Bearer "+token)

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a context-aware handler
    handler := http.HandlerFunc(middleware.GetTodos)
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to get todos")
    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var todos []models.Todo
    err = json.NewDecoder(rr.Body).Decode(&todos)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    fmt.Printf("Retrieved %d todos\n", len(todos))
    fmt.Println("GetTodos test passed")
}

// TestCreateTodo tests the CreateTodo endpoint
func TestCreateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTodo")
    fmt.Println("Testing creating a new todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    // Create a request body with todo details
    var jsonStr = []byte(`{"task":"Test Task", "description":"Test Description", "important":true}`)
    req, err := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a handler
    handler := http.HandlerFunc(middleware.CreateTodo)
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to create todo")
    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var todo models.Todo
    err = json.NewDecoder(rr.Body).Decode(&todo)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate the response
    if todo.Task != "Test Task" {
        t.Errorf("Handler returned unexpected task: got %v want %v", todo.Task, "Test Task")
    }
    if todo.Description != "Test Description" {
        t.Errorf("Handler returned unexpected description: got %v want %v", todo.Description, "Test Description")
    }
    if !todo.Important {
        t.Errorf("Handler returned unexpected important flag: got %v want %v", todo.Important, true)
    }
    if todo.Done {
        t.Errorf("Handler returned unexpected done flag: got %v want %v", todo.Done, false)
    }
    if todo.ID == "" {
        t.Errorf("Handler returned empty todo ID")
    } else {
        fmt.Printf("Todo created with ID: %s\n", todo.ID)
    }
    
    // Save the todo ID for update/delete tests
    todoID = todo.ID
    
    fmt.Println("CreateTodo test passed")
}

// TestUpdateTodo tests the UpdateTodo endpoint
func TestUpdateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateTodo")
    fmt.Println("Testing updating an existing todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    if todoID == "" {
        // Create a todo first if ID is not available
        TestCreateTodo(t)
    }
    
    // Create a request body with updated todo details
    var jsonStr = []byte(`{"task":"Updated Task", "description":"Updated Description", "important":true, "done":true}`)
    req, err := http.NewRequest("PUT", "/api/todo/"+todoID, bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Set up router with variables
    router := mux.NewRouter()
    router.HandleFunc("/api/todo/{id}", middleware.UpdateTodo).Methods("PUT")
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to update todo")
    router.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var todo models.Todo
    err = json.NewDecoder(rr.Body).Decode(&todo)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate the response
    if todo.Task != "Updated Task" {
        t.Errorf("Handler returned unexpected task: got %v want %v", todo.Task, "Updated Task")
    }
    if todo.Description != "Updated Description" {
        t.Errorf("Handler returned unexpected description: got %v want %v", todo.Description, "Updated Description")
    }
    if !todo.Done {
        t.Errorf("Handler returned unexpected done flag: got %v want %v", todo.Done, true)
    }
    
    fmt.Println("UpdateTodo test passed")
}

// TestUndoTodo tests the UndoTodo endpoint
func TestUndoTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUndoTodo")
    fmt.Println("Testing undoing a completed todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    if todoID == "" {
        // We need a todo that is marked as done
        TestUpdateTodo(t)
    }
    
    // Create a new request
    req, err := http.NewRequest("PUT", "/api/todo/undo/"+todoID, nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }

    // Set up router with variables
    router := mux.NewRouter()
    router.HandleFunc("/api/todo/undo/{id}", middleware.UndoTodo).Methods("PUT")
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to undo todo")
    router.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var todo models.Todo
    err = json.NewDecoder(rr.Body).Decode(&todo)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate the response - check that done is now false
    if todo.Done {
        t.Errorf("Handler returned unexpected done flag: got %v want %v", todo.Done, false)
    } else {
        fmt.Println("Todo was successfully marked as undone")
    }
    
    fmt.Println("UndoTodo test passed")
}

// TestDeleteTodo tests the DeleteTodo endpoint
func TestDeleteTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDeleteTodo")
    fmt.Println("Testing deleting a todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    if todoID == "" {
        // Create a todo first if ID is not available
        TestCreateTodo(t)
        if todoID == "" {
            t.Fatal("Failed to create a todo for deletion test")
        }
    }
    
    // Create a new request
    req, err := http.NewRequest("DELETE", "/api/todo/"+todoID, nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }

    // Set up router with variables
    router := mux.NewRouter()
    router.HandleFunc("/api/todo/{id}", middleware.DeleteTodo).Methods("DELETE")
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to delete todo")
    router.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var response map[string]string
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate the response
    if response["result"] != "success" {
        t.Errorf("Handler returned unexpected result: got %v want %v", response["result"], "success")
    } else {
        fmt.Println("Todo was successfully deleted")
    }
    
    fmt.Println("DeleteTodo test passed")
    
    // Clear the todoID since it's been deleted
    todoID = ""
    
    // Clean up after all todo tests
    cleanupTodoTestData()
}

// Global variable to store a todo ID for use in update/delete tests
var todoID string