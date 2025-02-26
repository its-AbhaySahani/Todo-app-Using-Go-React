package TestCases

import (
    "fmt"  // Add this import
    "os"
    "testing"
    "log"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
)

func TestMain(m *testing.M) {
    // Setup
    database.Connect()
    fmt.Println("🔄 Test database connected")
    
    // Run tests
    code := m.Run()
    
    // Clean up
    cleanupTestDataMain()  // Renamed to avoid conflict
	
	if code == 0 {
        fmt.Println("\n✅ ALL TESTS PASSED")
    } else {
        fmt.Println("\n❌ SOME TESTS FAILED")
    }
    
    
    // Exit
    os.Exit(code)
}

func setup() {
    // Connect to the database
    database.Connect()
    log.Println("Test database connected")
    
    // Clear test data to ensure clean state
    cleanupTestDataMain()  // Renamed to avoid conflict
}

func teardown() {
    // Clean up test data
    cleanupTestDataMain()  // Renamed to avoid conflict
    
    // Close database connection
    if database.DB != nil {
        err := database.DB.Close()
        if err != nil {
            log.Fatal("Error closing the database:", err)
        }
        log.Println("Test database connection closed")
    }
}

// Renamed to avoid conflict with the function in auth_test.go
func cleanupTestDataMain() {
    // Delete test user if exists
    _, err := database.DB.Exec("DELETE FROM users WHERE username = ?", "testuser")
    if err != nil {
        log.Println("Error cleaning up test data:", err)
    }
}