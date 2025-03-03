package TestCases

import (
    "fmt"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Test User ID that will be used consistently across all tests
const testUserID = "todo-test-user-id"
const testUsername = "todo_testuser" // Different from "testuser" used in auth tests

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

// Global variable to store a todo ID for use in update/delete tests
var todoID string

// TestGetTodos tests the GetTodos function directly
func TestGetTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTodos")
    fmt.Println("Testing getting todos for the user")
    
    // Clean up before test to ensure fresh state
    cleanupTodoTestData()
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    // First create a few test todos for the user
    task1 := "Test Task 1"
    desc1 := "Test Description 1"
    
    // Create first todo directly using the model function
    todo1, err := models.CreateTodo(task1, desc1, true, testUserID)
    if err != nil {
        t.Fatalf("Failed to create first test todo: %v", err)
    }
    fmt.Printf("Created first todo with ID: %s\n", todo1.ID)
    
    // Create second todo
    task2 := "Test Task 2"
    desc2 := "Test Description 2"
    todo2, err := models.CreateTodo(task2, desc2, false, testUserID)
    if err != nil {
        t.Fatalf("Failed to create second test todo: %v", err)
    }
    fmt.Printf("Created second todo with ID: %s\n", todo2.ID)
    
    // Now retrieve all todos for the user
    todos, err := models.GetTodos(testUserID)
    if err != nil {
        t.Fatalf("Failed to get todos: %v", err)
    }
    
    // Verify we retrieved the correct number of todos
    if len(todos) != 2 {
        t.Errorf("Expected 2 todos, got %d", len(todos))
    } else {
        fmt.Printf("Retrieved %d todos as expected\n", len(todos))
    }
    
    // Verify the todo details are correct
    foundTodo1 := false
    foundTodo2 := false
    
    for _, todo := range todos {
        if todo.ID == todo1.ID {
            foundTodo1 = true
            if todo.Task != task1 {
                t.Errorf("Expected task '%s', got '%s'", task1, todo.Task)
            }
            if todo.Description != desc1 {
                t.Errorf("Expected description '%s', got '%s'", desc1, todo.Description)
            }
            if !todo.Important {
                t.Errorf("Expected todo to be important, but it's not")
            }
        } else if todo.ID == todo2.ID {
            foundTodo2 = true
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
    
    if !foundTodo1 {
        t.Errorf("First todo not found in retrieved todos")
    }
    if !foundTodo2 {
        t.Errorf("Second todo not found in retrieved todos")
    }
    
    // Save one ID for later tests
    todoID = todo1.ID
    
    fmt.Println("GetTodos test passed")
}

// TestCreateTodo tests the CreateTodo function directly
func TestCreateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTodo")
    fmt.Println("Testing creating a new todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    // Create a new todo directly using the model function
    task := "Functional Test Task"
    description := "Functional Test Description"
    important := true
    
    todo, err := models.CreateTodo(task, description, important, testUserID)
    if err != nil {
        t.Fatalf("Failed to create todo: %v", err)
    }
    
    // Verify the todo was created with the correct data
    if todo.ID == "" {
        t.Errorf("Todo was created but ID is empty")
    } else {
        fmt.Printf("Created todo with ID: %s\n", todo.ID)
    }
    
    if todo.Task != task {
        t.Errorf("Expected task '%s', got '%s'", task, todo.Task)
    }
    
    if todo.Description != description {
        t.Errorf("Expected description '%s', got '%s'", description, todo.Description)
    }
    
    if !todo.Important {
        t.Errorf("Expected todo to be important, but it's not")
    }
    
    if todo.Done {
        t.Errorf("Expected todo to not be done, but it is")
    }
    
    // Verify the todo exists in the database
    var dbTask, dbDescription string
    var dbDone, dbImportant bool
    
    err = database.DB.QueryRow(
        "SELECT task, description, done, important FROM todos WHERE id = ? AND user_id = ?", 
        todo.ID, testUserID,
    ).Scan(&dbTask, &dbDescription, &dbDone, &dbImportant)
    
    if err != nil {
        t.Fatalf("Failed to retrieve todo from database: %v", err)
    }
    
    if dbTask != task {
        t.Errorf("Database task '%s' doesn't match expected task '%s'", dbTask, task)
    }
    
    if dbDescription != description {
        t.Errorf("Database description '%s' doesn't match expected description '%s'", dbDescription, description)
    }
    
    if !dbImportant {
        t.Errorf("Expected todo to be important in database, but it's not")
    }
    
    if dbDone {
        t.Errorf("Expected todo to not be done in database, but it is")
    }
    
    // Save ID for later tests if needed
    todoID = todo.ID
    
    fmt.Println("CreateTodo test passed")
}

// TestUpdateTodo tests the UpdateTodo function directly
func TestUpdateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateTodo")
    fmt.Println("Testing updating an existing todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    if todoID == "" {
        // Create a todo first if ID is not available
        TestCreateTodo(t)
    }
    
    // Update the todo directly using the model function
    updatedTask := "Updated Functional Task"
    updatedDesc := "Updated Functional Description"
    updatedDone := true
    updatedImportant := false
    
    updatedTodo, err := models.UpdateTodo(todoID, updatedTask, updatedDesc, updatedDone, updatedImportant, testUserID)
    if err != nil {
        t.Fatalf("Failed to update todo: %v", err)
    }
    
    // Verify the todo was updated with the correct data
    if updatedTodo.ID != todoID {
        t.Errorf("Expected todo ID '%s', got '%s'", todoID, updatedTodo.ID)
    }
    
    if updatedTodo.Task != updatedTask {
        t.Errorf("Expected task '%s', got '%s'", updatedTask, updatedTodo.Task)
    }
    
    if updatedTodo.Description != updatedDesc {
        t.Errorf("Expected description '%s', got '%s'", updatedDesc, updatedTodo.Description)
    }
    
    if !updatedTodo.Done {
        t.Errorf("Expected todo to be done, but it's not")
    }
    
    if updatedTodo.Important {
        t.Errorf("Expected todo to not be important, but it is")
    }
    
    // Verify the todo was updated in the database
    var dbTask, dbDescription string
    var dbDone, dbImportant bool
    
    err = database.DB.QueryRow(
        "SELECT task, description, done, important FROM todos WHERE id = ? AND user_id = ?", 
        todoID, testUserID,
    ).Scan(&dbTask, &dbDescription, &dbDone, &dbImportant)
    
    if err != nil {
        t.Fatalf("Failed to retrieve updated todo from database: %v", err)
    }
    
    if dbTask != updatedTask {
        t.Errorf("Database task '%s' doesn't match expected task '%s'", dbTask, updatedTask)
    }
    
    if dbDescription != updatedDesc {
        t.Errorf("Database description '%s' doesn't match expected description '%s'", dbDescription, updatedDesc)
    }
    
    if !dbDone {
        t.Errorf("Expected todo to be done in database, but it's not")
    }
    
    if dbImportant {
        t.Errorf("Expected todo to not be important in database, but it is")
    }
    
    fmt.Println("UpdateTodo test passed")
}

// TestUndoTodo tests the UndoTodo function directly
func TestUndoTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUndoTodo")
    fmt.Println("Testing undoing a completed todo")
    
    // Ensure test user exists
    ensureTestUserExists(t)
    
    if todoID == "" {
        // We need a todo that is marked as done
        TestUpdateTodo(t)
    }
    
    // Ensure the todo is marked as done before testing undo
    var isDone bool
    err := database.DB.QueryRow(
        "SELECT done FROM todos WHERE id = ? AND user_id = ?", 
        todoID, testUserID,
    ).Scan(&isDone)
    
    if err != nil {
        t.Fatalf("Failed to check if todo is done: %v", err)
    }
    
    if !isDone {
        // If not done, update it to be done
        _, err := database.DB.Exec(
            "UPDATE todos SET done = ? WHERE id = ? AND user_id = ?", 
            true, todoID, testUserID,
        )
        
        if err != nil {
            t.Fatalf("Failed to mark todo as done for undo test: %v", err)
        }
    }
    
    // Now undo the todo directly using the model function
    undoneTask, err := models.UndoTodo(todoID, testUserID)
    if err != nil {
        t.Fatalf("Failed to undo todo: %v", err)
    }
    
    // Verify the todo was undone
    if undoneTask.ID != todoID {
        t.Errorf("Expected todo ID '%s', got '%s'", todoID, undoneTask.ID)
    }
    
    if undoneTask.Done {
        t.Errorf("Expected todo to be undone (not done), but it's still marked as done")
    }
    
    // Verify the todo was undone in the database
    var dbDone bool
    
    err = database.DB.QueryRow(
        "SELECT done FROM todos WHERE id = ? AND user_id = ?", 
        todoID, testUserID,
    ).Scan(&dbDone)
    
    if err != nil {
        t.Fatalf("Failed to retrieve undone todo from database: %v", err)
    }
    
    if dbDone {
        t.Errorf("Expected todo to be undone (not done) in database, but it's still marked as done")
    } else {
        fmt.Println("Todo was successfully marked as undone in database")
    }
    
    fmt.Println("UndoTodo test passed")
}

// TestDeleteTodo tests the DeleteTodo function directly
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
    
    // Delete the todo directly using the model function
    err := models.DeleteTodo(todoID, testUserID)
    if err != nil {
        t.Fatalf("Failed to delete todo: %v", err)
    }
    
    // Verify the todo was deleted from the database
    var count int
    err = database.DB.QueryRow(
        "SELECT COUNT(*) FROM todos WHERE id = ? AND user_id = ?", 
        todoID, testUserID,
    ).Scan(&count)
    
    if err != nil {
        t.Fatalf("Failed to check if todo was deleted: %v", err)
    }
    
    if count != 0 {
        t.Errorf("Expected todo to be deleted from database, but it still exists")
    } else {
        fmt.Println("Todo was successfully deleted from database")
    }
    
    fmt.Println("DeleteTodo test passed")
    
    // Clear the todoID since it's been deleted
    todoID = ""
    
    // Clean up after all todo tests
    cleanupTodoTestData()
}

// TestGetDailyRoutines tests the GetDailyRoutines function
func TestGetDailyRoutines(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetDailyRoutines")
    fmt.Println("Testing retrieving routines for a specific day and schedule type")
    
    // Ensure test user exists and clean up data
    cleanupTodoTestData()
    ensureTestUserExists(t)
    
    // Create a test todo
    todo, err := models.CreateTodo("Routine Test Task", "Test routine task", false, testUserID)
    if err != nil {
        t.Fatalf("Failed to create todo for routine test: %v", err)
    }
    
    // Get today's day name in lowercase
    today := time.Now().Weekday().String()
    day := today
    scheduleType := "morning"
    
    // Create a routine directly in the database
    routineID := uuid.New().String()
    currentDate := time.Now().Format("2006-01-02")
    
    _, err = database.DB.Exec(
        "INSERT INTO routines (id, day, scheduleType, taskId, userId, createdAt, updatedAt, isActive) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        routineID, day, scheduleType, todo.ID, testUserID, currentDate, currentDate, true,
    )
    
    if err != nil {
        t.Fatalf("Failed to create routine: %v", err)
    }
    
    // Now test the GetDailyRoutines function
    routineTodos, err := models.GetDailyRoutines(day, scheduleType, testUserID)
    if err != nil {
        t.Fatalf("Failed to get daily routines: %v", err)
    }
    
    // Verify we got the correct routines
    if len(routineTodos) != 1 {
        t.Errorf("Expected 1 routine todo, got %d", len(routineTodos))
    } else {
        fmt.Printf("Retrieved %d routine todos as expected\n", len(routineTodos))
        
        // Verify the routine todo details
        routineTodo := routineTodos[0]
        if routineTodo.ID != todo.ID {
            t.Errorf("Expected todo ID '%s', got '%s'", todo.ID, routineTodo.ID)
        }
        
        if routineTodo.Task != "Routine Test Task" {
            t.Errorf("Expected task 'Routine Test Task', got '%s'", routineTodo.Task)
        }
    }
    
    // Clean up by deleting the routine
    _, err = database.DB.Exec("DELETE FROM routines WHERE id = ?", routineID)
    if err != nil {
        fmt.Printf("Warning: Failed to delete test routine: %v\n", err)
    }
    
    // Clean up the todo
    cleanupTodoTestData()
    
    fmt.Println("GetDailyRoutines test passed")
}

// TestUpdateRoutineDay tests the UpdateRoutineDay function
func TestUpdateRoutineDay(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateRoutineDay")
    fmt.Println("Testing updating a routine's day")
    
    // Ensure test user exists and clean up data
    cleanupTodoTestData()
    ensureTestUserExists(t)
    
    // Create a test todo
    todo, err := models.CreateTodo("Routine Day Update Test", "Test routine day update", false, testUserID)
    if err != nil {
        t.Fatalf("Failed to create todo for routine test: %v", err)
    }
    
    // Create a routine directly in the database
    routineID := uuid.New().String()
    currentDate := time.Now().Format("2006-01-02")
    initialDay := "monday"
    scheduleType := "evening"
    
    _, err = database.DB.Exec(
        "INSERT INTO routines (id, day, scheduleType, taskId, userId, createdAt, updatedAt, isActive) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        routineID, initialDay, scheduleType, todo.ID, testUserID, currentDate, currentDate, true,
    )
    
    if err != nil {
        t.Fatalf("Failed to create routine: %v", err)
    }
    
    // Update the routine day
    newDay := "friday"
    err = models.UpdateRoutineDay(routineID, newDay)
    if err != nil {
        t.Fatalf("Failed to update routine day: %v", err)
    }
    
    // Verify the day was updated in the database
    var dbDay string
    err = database.DB.QueryRow(
        "SELECT day FROM routines WHERE id = ?", 
        routineID,
    ).Scan(&dbDay)
    
    if err != nil {
        t.Fatalf("Failed to retrieve updated routine from database: %v", err)
    }
    
    if dbDay != newDay {
        t.Errorf("Expected routine day '%s', got '%s'", newDay, dbDay)
    } else {
        fmt.Println("Routine day was successfully updated in database")
    }
    
    // Clean up by deleting the routine
    _, err = database.DB.Exec("DELETE FROM routines WHERE id = ?", routineID)
    if err != nil {
        fmt.Printf("Warning: Failed to delete test routine: %v\n", err)
    }
    
    // Clean up the todo
    cleanupTodoTestData()
    
    fmt.Println("UpdateRoutineDay test passed")
}
