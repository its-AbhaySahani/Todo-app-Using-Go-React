package TestCases

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Test IDs that will be used consistently across all tests
const sharedRecipientUserID = "shared-recipient-user-id"
const sharedRecipientUsername = "shared_recipient_user"

// Helper function to ensure the recipient test user exists
func ensureRecipientUserExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", sharedRecipientUserID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if recipient user exists: %v", err)
    }

    if count == 0 {
        // Create the recipient user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            sharedRecipientUserID, sharedRecipientUsername, "$2a$10$TestHashedPasswordForRecipientUser",
        )
        if err != nil {
            t.Fatalf("Failed to create recipient user: %v", err)
        }
        fmt.Println("Created recipient user with ID:", sharedRecipientUserID)
    } else {
        fmt.Println("Recipient user already exists with ID:", sharedRecipientUserID)
    }
}

// Helper function to cleanup shared test data
func cleanupSharedTestData() {
    // Delete shared todos
    _, err := database.DB.Exec("DELETE FROM shared_todos WHERE user_id = ? OR shared_by = ?", 
        sharedRecipientUserID, testUserID)
    if err != nil {
        fmt.Println("Error cleaning up shared todos:", err)
    } else {
        fmt.Println("Shared todos cleaned up")
    }
    
    // Delete recipient user
    _, err = database.DB.Exec("DELETE FROM users WHERE id = ?", sharedRecipientUserID)
    if err != nil {
        fmt.Println("Error cleaning up recipient user:", err)
    } else {
        fmt.Println("Recipient user cleaned up")
    }
}

// TestShareTodo tests the ShareTodo endpoint
func TestShareTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestShareTodo")
    fmt.Println("Testing sharing a todo with another user")
    
    // Clean up before test to ensure fresh state
    cleanupSharedTestData()
    
    // Ensure both users exist
    ensureTestUserExists(t)
    ensureRecipientUserExists(t)
    
    // First, create a todo to share
    var sharedTodoID string
    
    // Create a todo for the sender user
    var jsonStr = []byte(`{"task":"Shared Task", "description":"Task to be shared", "important":true}`)
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
    
    // Serve the request to create todo
    fmt.Println("Creating a todo to be shared")
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }
    
    // Parse the response to get the todo ID
    var todo models.Todo
    err = json.NewDecoder(rr.Body).Decode(&todo)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }
    
    sharedTodoID = todo.ID
    fmt.Printf("Created todo with ID: %s to be shared\n", sharedTodoID)
    
    // Now share the todo with the recipient user
    shareJSON := []byte(`{"taskId":"` + sharedTodoID + `", "username":"` + sharedRecipientUsername + `"}`)
    req, err = http.NewRequest("POST", "/api/share", bytes.NewBuffer(shareJSON))
    if err != nil {
        t.Fatal("Failed to create share request:", err)
    }
    req.Header.Set("Content-Type", "application/json")
    
    // Create a new response recorder
    rr = httptest.NewRecorder()
    
    // Create the share handler
    shareHandler := http.HandlerFunc(middleware.ShareTodo)
    
    // Create a request with user ID in context
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request to share todo
    fmt.Println("Sending request to share todo")
    shareHandler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Share handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }
    
    // Parse the response
    var response map[string]string
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to parse share response body: %v\nResponse body: %s", err, rr.Body.String())
    }
    
    // Validate the response
    if response["result"] != "success" {
        t.Errorf("Handler returned unexpected result: got %v want %v", response["result"], "success")
    } else {
        fmt.Println("Todo was successfully shared")
    }
    
    fmt.Println("ShareTodo test passed")
}

// TestGetSharedTodos tests the GetSharedTodos endpoint
func TestGetSharedTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetSharedTodos")
    fmt.Println("Testing getting shared todos")
    
    // First share a todo to ensure there's data to retrieve
    TestShareTodo(t)
    
    // Create a new request
    req, err := http.NewRequest("GET", "/api/shared", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create handler
    handler := http.HandlerFunc(middleware.GetSharedTodos)
    
    // Create a request with recipient user ID in context to check received todos
    req = addUserIDToContext(req, sharedRecipientUserID)
    
    // Serve the request
    fmt.Println("Sending request to get shared todos (as recipient)")
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }
    
    // Parse the response
    var responseRecipient map[string]interface{}
    err = json.NewDecoder(rr.Body).Decode(&responseRecipient)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }
    
    // Check that we got received todos
    received, ok := responseRecipient["received"].([]interface{})
    if !ok {
        t.Errorf("Response missing 'received' key or it's not an array")
    } else {
        fmt.Printf("Recipient has %d received shared todos\n", len(received))
        if len(received) == 0 {
            t.Errorf("Expected at least one received todo")
        }
    }
    
    // Now check the shared todos from the sender's perspective
    req, err = http.NewRequest("GET", "/api/shared", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    
    // Create a new response recorder
    rr = httptest.NewRecorder()
    
    // Create a request with sender user ID in context to check shared todos
    req = addUserIDToContext(req, testUserID)
    
    // Serve the request
    fmt.Println("Sending request to get shared todos (as sender)")
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }
    
    // Parse the response
    var responseSender map[string]interface{}
    err = json.NewDecoder(rr.Body).Decode(&responseSender)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }
    
    // Check that we got shared todos
    shared, ok := responseSender["shared"].([]interface{})
    if !ok {
        t.Errorf("Response missing 'shared' key or it's not an array")
    } else {
        fmt.Printf("Sender has shared %d todos\n", len(shared))
        if len(shared) == 0 {
            t.Errorf("Expected at least one shared todo")
        }
    }
    
    fmt.Println("GetSharedTodos test passed")
    
    // Clean up shared test data
    cleanupSharedTestData()
}