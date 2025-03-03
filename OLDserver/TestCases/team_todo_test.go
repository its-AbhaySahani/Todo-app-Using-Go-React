package TestCases

import (
    "fmt"
    "testing"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Constants for team todo tests
const testTeamTodoUserID = "team-todo-test-user-id"
const testTeamTodoUsername = "team_todo_testuser"
var testTeamTodoID string // Will be populated when we create a test team todo
var testTeamForTodoID string // Will be populated when we create a test team

// Helper function to ensure the test user exists
func ensureTeamTodoTestUserExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testTeamTodoUserID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if team todo test user exists: %v", err)
    }

    if count == 0 {
        // Create the test user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testTeamTodoUserID, testTeamTodoUsername, "$2a$10$TestHashedPasswordForTeamTodoTestUser",
        )
        if err != nil {
            t.Fatalf("Failed to create team todo test user: %v", err)
        }
        fmt.Println("Created team todo test user with ID:", testTeamTodoUserID)
    } else {
        fmt.Println("Team todo test user already exists with ID:", testTeamTodoUserID)
    }
}

// Helper function to ensure test team exists
func ensureTestTeamForTodoExists(t *testing.T) string {
    // If we already have a team ID, check if it still exists
    if testTeamForTodoID != "" {
        var count int
        err := database.DB.QueryRow("SELECT COUNT(*) FROM teams WHERE id = ?", testTeamForTodoID).Scan(&count)
        if err == nil && count > 0 {
            fmt.Println("Using existing team with ID:", testTeamForTodoID)
            return testTeamForTodoID
        }
    }
    
    // Create a new team
    teamID := uuid.New().String()
    teamName := fmt.Sprintf("Test Team %s", uuid.New().String()[:8])
    teamPassword := "testteampassword"
    
    // Make sure the test user exists
    ensureTeamTodoTestUserExists(t)
    
    _, err := database.DB.Exec(
        "INSERT INTO teams (id, name, password, admin_id) VALUES (?, ?, ?, ?)",
        teamID, teamName, teamPassword, testTeamTodoUserID,
    )
    
    if err != nil {
        t.Fatalf("Failed to create test team: %v", err)
    }
    
    // Add the user as a team member with admin privileges
    _, err = database.DB.Exec(
        "INSERT INTO team_members (team_id, user_id, is_admin) VALUES (?, ?, ?)",
        teamID, testTeamTodoUserID, true,
    )
    
    if err != nil {
        t.Fatalf("Failed to add user as team member: %v", err)
    }
    
    fmt.Println("Created test team with ID:", teamID)
    testTeamForTodoID = teamID
    return teamID
}

// Helper function to cleanup team test data after tests
func cleanupTeamTodoTestData() {
    // Delete team todos
    if testTeamForTodoID != "" {
        _, err := database.DB.Exec("DELETE FROM team_todos WHERE team_id = ?", testTeamForTodoID)
        if err != nil {
            fmt.Println("Error cleaning up team todos:", err)
        } else {
            fmt.Println("Team todos cleaned up")
        }
        
        // Delete team members
        _, err = database.DB.Exec("DELETE FROM team_members WHERE team_id = ?", testTeamForTodoID)
        if err != nil {
            fmt.Println("Error cleaning up team members:", err)
        } else {
            fmt.Println("Team members cleaned up")
        }
        
        // Delete the team
        _, err = database.DB.Exec("DELETE FROM teams WHERE id = ?", testTeamForTodoID)
        if err != nil {
            fmt.Println("Error cleaning up team:", err)
        } else {
            fmt.Println("Team cleaned up")
        }
    }
    
    // Delete test user
    _, err := database.DB.Exec("DELETE FROM users WHERE id = ?", testTeamTodoUserID)
    if err != nil {
        fmt.Println("Error cleaning up team test user:", err)
    } else {
        fmt.Println("Team test user cleaned up")
    }
    
    testTeamForTodoID = "" // Reset the team ID
    testTeamTodoID = ""    // Reset the team todo ID
}

// TestCreateTeamTodo tests the CreateTeamTodo function directly
func TestCreateTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTeamTodo")
    fmt.Println("Testing creating a new team todo")
    
    // Clean up before test to ensure fresh state
    cleanupTeamTodoTestData()
    
    // Ensure test team exists
    teamID := ensureTestTeamForTodoExists(t)
    
    // Create a new team todo directly using the model function
    task := "Team Functional Test Task"
    description := "Team Functional Test Description"
    important := true
    assignedTo := testTeamTodoUserID
    
    teamTodo, err := models.CreateTeamTodo(task, description, important, teamID, assignedTo)
    if err != nil {
        t.Fatalf("Failed to create team todo: %v", err)
    }
    
    // Verify the team todo was created with the correct data
    if teamTodo.ID == "" {
        t.Errorf("Team todo was created but ID is empty")
    } else {
        fmt.Printf("Created team todo with ID: %s\n", teamTodo.ID)
    }
    
    if teamTodo.Task != task {
        t.Errorf("Expected task '%s', got '%s'", task, teamTodo.Task)
    }
    
    if teamTodo.Description != description {
        t.Errorf("Expected description '%s', got '%s'", description, teamTodo.Description)
    }
    
    if !teamTodo.Important {
        t.Errorf("Expected team todo to be important, but it's not")
    }
    
    if teamTodo.Done {
        t.Errorf("Expected team todo to not be done, but it is")
    }
    
    if teamTodo.TeamID != teamID {
        t.Errorf("Expected team ID '%s', got '%s'", teamID, teamTodo.TeamID)
    }
    
    if teamTodo.AssignedTo != assignedTo {
        t.Errorf("Expected assigned to '%s', got '%s'", assignedTo, teamTodo.AssignedTo)
    }
    
    // Verify the team todo exists in the database
    var dbTask, dbDescription, dbTeamID, dbAssignedTo string
    var dbDone, dbImportant bool
    
    err = database.DB.QueryRow(
        "SELECT task, description, done, important, team_id, assigned_to FROM team_todos WHERE id = ?", 
        teamTodo.ID,
    ).Scan(&dbTask, &dbDescription, &dbDone, &dbImportant, &dbTeamID, &dbAssignedTo)
    
    if err != nil {
        t.Fatalf("Failed to retrieve team todo from database: %v", err)
    }
    
    if dbTask != task {
        t.Errorf("Database task '%s' doesn't match expected task '%s'", dbTask, task)
    }
    
    if dbDescription != description {
        t.Errorf("Database description '%s' doesn't match expected description '%s'", dbDescription, description)
    }
    
    if !dbImportant {
        t.Errorf("Expected team todo to be important in database, but it's not")
    }
    
    if dbDone {
        t.Errorf("Expected team todo to not be done in database, but it is")
    }
    
    if dbTeamID != teamID {
        t.Errorf("Expected team ID '%s' in database, got '%s'", teamID, dbTeamID)
    }
    
    if dbAssignedTo != assignedTo {
        t.Errorf("Expected assigned to '%s' in database, got '%s'", assignedTo, dbAssignedTo)
    }
    
    // Save ID for later tests
    testTeamTodoID = teamTodo.ID
    
    fmt.Println("CreateTeamTodo test passed")
}

// TestGetTeamTodos tests the GetTeamTodos function directly
func TestGetTeamTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamTodos")
    fmt.Println("Testing getting todos for the team")
    
    // Ensure test team exists
    teamID := ensureTestTeamForTodoExists(t)
    
    if testTeamTodoID == "" {
        // Create a todo first if ID is not available
        TestCreateTeamTodo(t)
        if testTeamTodoID == "" {
            t.Fatal("Failed to create a team todo for retrieval test")
        }
    }
    
    // Create a second team todo for testing retrieval of multiple todos
    task2 := "Team Task 2"
    desc2 := "Team Description 2"
    
    todo2, err := models.CreateTeamTodo(task2, desc2, false, teamID, testTeamTodoUserID)
    if err != nil {
        t.Fatalf("Failed to create second test team todo: %v", err)
    }
    fmt.Printf("Created second team todo with ID: %s\n", todo2.ID)
    
    // Now retrieve all todos for the team
    todos, err := models.GetTeamTodos(teamID)
    if err != nil {
        t.Fatalf("Failed to get team todos: %v", err)
    }
    
    // Verify we retrieved the correct number of todos (should be at least 2)
    if len(todos) < 2 {
        t.Errorf("Expected at least 2 team todos, got %d", len(todos))
    } else {
        fmt.Printf("Retrieved %d team todos\n", len(todos))
    }
    
    // Verify we can find both todos in the results
    found1 := false
    found2 := false
    
    for _, todo := range todos {
        if todo.ID == testTeamTodoID {
            found1 = true
        } else if todo.ID == todo2.ID {
            found2 = true
            if todo.Task != task2 {
                t.Errorf("Expected task '%s', got '%s'", task2, todo.Task)
            }
            if todo.Description != desc2 {
                t.Errorf("Expected description '%s', got '%s'", desc2, todo.Description)
            }
            if todo.Important {
                t.Errorf("Expected todo to not be important, but it is")
            }
        }
    }
    
    if !found1 {
        t.Errorf("First team todo not found in retrieved todos")
    }
    if !found2 {
        t.Errorf("Second team todo not found in retrieved todos")
    }
    
    fmt.Println("GetTeamTodos test passed")
}

// TestUpdateTeamTodo tests the UpdateTeamTodo function directly
func TestUpdateTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateTeamTodo")
    fmt.Println("Testing updating an existing team todo")
    
    // Ensure test team exists
    teamID := ensureTestTeamForTodoExists(t)
    
    if testTeamTodoID == "" {
        // Create a team todo first if ID is not available
        TestCreateTeamTodo(t)
        if testTeamTodoID == "" {
            t.Fatal("Failed to create a team todo for update test")
        }
    }
    
    // Update the team todo directly using the model function
    updatedTask := "Updated Team Task"
    updatedDesc := "Updated Team Description"
    updatedDone := true
    updatedImportant := false
    updatedAssignedTo := testTeamTodoUserID // Keep the same assigned user for simplicity
    
    updatedTodo, err := models.UpdateTeamTodo(testTeamTodoID, updatedTask, updatedDesc, updatedDone, updatedImportant, teamID, updatedAssignedTo)
    if err != nil {
        t.Fatalf("Failed to update team todo: %v", err)
    }
    
    // Verify the team todo was updated with the correct data
    if updatedTodo.ID != testTeamTodoID {
        t.Errorf("Expected team todo ID '%s', got '%s'", testTeamTodoID, updatedTodo.ID)
    }
    
    if updatedTodo.Task != updatedTask {
        t.Errorf("Expected task '%s', got '%s'", updatedTask, updatedTodo.Task)
    }
    
    if updatedTodo.Description != updatedDesc {
        t.Errorf("Expected description '%s', got '%s'", updatedDesc, updatedTodo.Description)
    }
    
    if !updatedTodo.Done {
        t.Errorf("Expected team todo to be done, but it's not")
    }
    
    if updatedTodo.Important {
        t.Errorf("Expected team todo to not be important, but it is")
    }
    
    // Verify the team todo was updated in the database
    var dbTask, dbDescription, dbAssignedTo string
    var dbDone, dbImportant bool
    
    err = database.DB.QueryRow(
        "SELECT task, description, done, important, assigned_to FROM team_todos WHERE id = ? AND team_id = ?", 
        testTeamTodoID, teamID,
    ).Scan(&dbTask, &dbDescription, &dbDone, &dbImportant, &dbAssignedTo)
    
    if err != nil {
        t.Fatalf("Failed to retrieve updated team todo from database: %v", err)
    }
    
    if dbTask != updatedTask {
        t.Errorf("Database task '%s' doesn't match expected task '%s'", dbTask, updatedTask)
    }
    
    if dbDescription != updatedDesc {
        t.Errorf("Database description '%s' doesn't match expected description '%s'", dbDescription, updatedDesc)
    }
    
    if !dbDone {
        t.Errorf("Expected team todo to be done in database, but it's not")
    }
    
    if dbImportant {
        t.Errorf("Expected team todo to not be important in database, but it is")
    }
    
    if dbAssignedTo != updatedAssignedTo {
        t.Errorf("Expected assigned to '%s' in database, got '%s'", updatedAssignedTo, dbAssignedTo)
    }
    
    fmt.Println("UpdateTeamTodo test passed")
}

// TestDeleteTeamTodo tests the DeleteTeamTodo function directly
func TestDeleteTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDeleteTeamTodo")
    fmt.Println("Testing deleting a team todo")
    
    // Ensure test team exists
    teamID := ensureTestTeamForTodoExists(t)
    
    if testTeamTodoID == "" {
        // Create a team todo first if ID is not available
        TestCreateTeamTodo(t)
        if testTeamTodoID == "" {
            t.Fatal("Failed to create a team todo for deletion test")
        }
    }
    
    // Delete the team todo directly using the model function
    err := models.DeleteTeamTodo(testTeamTodoID, teamID)
    if err != nil {
        t.Fatalf("Failed to delete team todo: %v", err)
    }
    
    // Verify the team todo was deleted from the database
    var count int
    err = database.DB.QueryRow(
        "SELECT COUNT(*) FROM team_todos WHERE id = ? AND team_id = ?", 
        testTeamTodoID, teamID,
    ).Scan(&count)
    
    if err != nil {
        t.Fatalf("Failed to check if team todo was deleted: %v", err)
    }
    
    if count != 0 {
        t.Errorf("Expected team todo to be deleted from database, but it still exists")
    } else {
        fmt.Println("Team todo was successfully deleted from database")
    }
    
    fmt.Println("DeleteTeamTodo test passed")
    
    // Clear the teamTodoID since it's been deleted
    testTeamTodoID = ""
    
    // Clean up after all team todo tests
    cleanupTeamTodoTestData()
}
