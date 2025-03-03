package TestCases

import (
    "fmt"
    "testing"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Test IDs that will be used consistently across all tests
const testSenderUserID = "shared-sender-user-id"
const testSenderUsername = "shared_sender_user"
const testRecipientUserID = "shared-recipient-user-id"
const testRecipientUsername = "shared_recipient_user"
var sharedTodoID string // Will be populated when we create a test todo

// Helper function to ensure the sender test user exists
func ensureSenderTestUserExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testSenderUserID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if sender user exists: %v", err)
    }

    if count == 0 {
        // Create the sender user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testSenderUserID, testSenderUsername, "$2a$10$TestHashedPasswordForSenderUser",
        )
        if err != nil {
            t.Fatalf("Failed to create sender user: %v", err)
        }
        fmt.Println("Created sender user with ID:", testSenderUserID)
    } else {
        fmt.Println("Sender user already exists with ID:", testSenderUserID)
    }
}

// Helper function to ensure the recipient test user exists
func ensureRecipientTestUserExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testRecipientUserID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if recipient user exists: %v", err)
    }

    if count == 0 {
        // Create the recipient user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testRecipientUserID, testRecipientUsername, "$2a$10$TestHashedPasswordForRecipientUser",
        )
        if err != nil {
            t.Fatalf("Failed to create recipient user: %v", err)
        }
        fmt.Println("Created recipient user with ID:", testRecipientUserID)
    } else {
        fmt.Println("Recipient user already exists with ID:", testRecipientUserID)
    }
}

// Helper function to cleanup shared test data
func cleanupSharedTestData() {
    // Delete shared todos
    _, err := database.DB.Exec("DELETE FROM shared_todos WHERE user_id = ? OR shared_by = ?", 
        testRecipientUserID, testSenderUserID)
    if err != nil {
        fmt.Println("Error cleaning up shared todos:", err)
    } else {
        fmt.Println("Shared todos cleaned up")
    }
    
    // Delete test todos
    _, err = database.DB.Exec("DELETE FROM todos WHERE user_id = ?", testSenderUserID)
    if err != nil {
        fmt.Println("Error cleaning up sender todos:", err)
    } else {
        fmt.Println("Sender todos cleaned up")
    }
    
    // Delete test users
    _, err = database.DB.Exec("DELETE FROM users WHERE id IN (?, ?)", testSenderUserID, testRecipientUserID)
    if err != nil {
        fmt.Println("Error cleaning up test users:", err)
    } else {
        fmt.Println("Test users cleaned up")
    }
}

// TestShareTodo tests the ShareTodoWithUser function directly
func TestShareTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestShareTodo")
    fmt.Println("Testing sharing a todo with another user")
    
    // Clean up before test to ensure fresh state
    cleanupSharedTestData()
    
    // Ensure both users exist
    ensureSenderTestUserExists(t)
    ensureRecipientTestUserExists(t)
    
    // First, create a todo to share
    task := "Shared Task"
    description := "Task to be shared"
    important := true
    
    // Create a new todo directly using the model function
    todo, err := models.CreateTodo(task, description, important, testSenderUserID)
    if err != nil {
        t.Fatalf("Failed to create todo: %v", err)
    }
    
    fmt.Printf("Created todo with ID: %s to be shared\n", todo.ID)
    sharedTodoID = todo.ID
    
    // Now share the todo with the recipient user
    err = models.ShareTodoWithUser(sharedTodoID, testRecipientUserID, testSenderUserID)
    if err != nil {
        t.Fatalf("Failed to share todo: %v", err)
    }
    
    // Verify the todo was shared by checking the shared_todos table
    var dbTask, dbDescription, dbUserID, dbSharedBy string
    var dbDone, dbImportant bool
    
    err = database.DB.QueryRow(
        "SELECT task, description, done, important, user_id, shared_by FROM shared_todos WHERE id = ?", 
        sharedTodoID,
    ).Scan(&dbTask, &dbDescription, &dbDone, &dbImportant, &dbUserID, &dbSharedBy)
    
    if err != nil {
        t.Fatalf("Failed to retrieve shared todo from database: %v", err)
    }
    
    if dbTask != task {
        t.Errorf("Database task '%s' doesn't match expected task '%s'", dbTask, task)
    }
    
    if dbDescription != description {
        t.Errorf("Database description '%s' doesn't match expected description '%s'", dbDescription, description)
    }
    
    if dbUserID != testRecipientUserID {
        t.Errorf("Expected user ID '%s', got '%s'", testRecipientUserID, dbUserID)
    }
    
    if dbSharedBy != testSenderUserID {
        t.Errorf("Expected shared by '%s', got '%s'", testSenderUserID, dbSharedBy)
    }
    
    if !dbImportant {
        t.Errorf("Expected shared todo to be important, but it's not")
    }
    
    fmt.Println("ShareTodo test passed")
}

// TestGetSharedTodos tests the GetSharedTodos function directly
func TestGetSharedTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetSharedTodos")
    fmt.Println("Testing retrieving shared todos for a user")
    
    // First share a todo to ensure there's data to retrieve
    if sharedTodoID == "" {
        TestShareTodo(t)
        if sharedTodoID == "" {
            t.Fatal("Failed to create and share a todo for retrieval test")
        }
    }
    
    // Get shared todos for the recipient
    todos, err := models.GetSharedTodos(testRecipientUserID)
    if err != nil {
        t.Fatalf("Failed to get shared todos: %v", err)
    }
    
    // Verify we got the correct shared todo
    if len(todos) == 0 {
        t.Errorf("Expected at least one shared todo, got none")
    } else {
        fmt.Printf("Retrieved %d shared todos\n", len(todos))
        
        // Find our shared todo
        found := false
        for _, todo := range todos {
            if todo.ID == sharedTodoID {
                found = true
                if todo.Task != "Shared Task" {
                    t.Errorf("Expected task 'Shared Task', got '%s'", todo.Task)
                }
                
                if todo.SharedBy != testSenderUserID {
                    t.Errorf("Expected shared by '%s', got '%s'", testSenderUserID, todo.SharedBy)
                }
                
                if todo.UserID != testRecipientUserID {
                    t.Errorf("Expected user ID '%s', got '%s'", testRecipientUserID, todo.UserID)
                }
            }
        }
        
        if !found {
            t.Errorf("Shared todo not found in retrieved todos")
        }
    }
    
    fmt.Println("GetSharedTodos test passed")
}

// TestGetSharedByMeTodos tests the GetSharedByMeTodos function directly
func TestGetSharedByMeTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetSharedByMeTodos")
    fmt.Println("Testing retrieving todos shared by a user")
    
    // First share a todo to ensure there's data to retrieve
    if sharedTodoID == "" {
        TestShareTodo(t)
        if sharedTodoID == "" {
            t.Fatal("Failed to create and share a todo for retrieval test")
        }
    }
    
    // Get todos shared by the sender
    todos, err := models.GetSharedByMeTodos(testSenderUserID)
    if err != nil {
        t.Fatalf("Failed to get shared by me todos: %v", err)
    }
    
    // Verify we got the correct shared todo
    if len(todos) == 0 {
        t.Errorf("Expected at least one shared by me todo, got none")
    } else {
        fmt.Printf("Retrieved %d shared by me todos\n", len(todos))
        
        // Find our shared todo
        found := false
        for _, todo := range todos {
            if todo.ID == sharedTodoID {
                found = true
                if todo.Task != "Shared Task" {
                    t.Errorf("Expected task 'Shared Task', got '%s'", todo.Task)
                }
                
                if todo.SharedBy != testSenderUserID {
                    t.Errorf("Expected shared by '%s', got '%s'", testSenderUserID, todo.SharedBy)
                }
                
                if todo.UserID != testRecipientUserID {
                    t.Errorf("Expected user ID '%s', got '%s'", testRecipientUserID, todo.UserID)
                }
            }
        }
        
        if !found {
            t.Errorf("Shared todo not found in shared by me todos")
        }
    }
    
    fmt.Println("GetSharedByMeTodos test passed")
    
    // Clean up shared test data
    cleanupSharedTestData()
}
