package unit

import (
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTodos")
    fmt.Println("Testing getting todos for the user")
    
    // Create a mock todo repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create test data
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    todos := []models.Todo{
        {
            ID:          "todo-1",
            Task:        "Test Task 1",
            Description: "Test Description 1",
            Done:        false,
            Important:   true,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        timeStr,
        },
        {
            ID:          "todo-2",
            Task:        "Test Task 2",
            Description: "Test Description 2",
            Done:        true,
            Important:   false,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        timeStr,
        },
    }
    
    // Setup expectations
    mockRepo.On("GetTodos", "user-123").Return(todos, nil)
    mockRepo.On("GetTodos", "empty-user").Return([]models.Todo{}, nil)
    mockRepo.On("GetTodos", "error-user").Return([]models.Todo{}, errors.New("database error"))
    
    // Test 1: User with todos
    userTodos, err := mockRepo.GetTodos("user-123")
    
    assert.NoError(t, err)
    assert.Equal(t, 2, len(userTodos))
    assert.Equal(t, "Test Task 1", userTodos[0].Task)
    assert.Equal(t, "Test Description 1", userTodos[0].Description)
    assert.True(t, userTodos[0].Important)
    assert.False(t, userTodos[0].Done)
    assert.Equal(t, "user-123", userTodos[0].UserID)
    
    fmt.Printf("Retrieved %d todos as expected\n", len(userTodos))
    for i, todo := range userTodos {
        fmt.Printf("Todo %d: {ID:%s Task:%s Description:%s Done:%t Important:%t}\n", 
            i+1, todo.ID, todo.Task, todo.Description, todo.Done, todo.Important)
    }
    
    // Test 2: User with no todos
    emptyTodos, err := mockRepo.GetTodos("empty-user")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyTodos))
    
    fmt.Println("Successfully retrieved empty todo list for user with no todos")
    
    // Test 3: Database error
    errorTodos, err := mockRepo.GetTodos("error-user")
    assert.Error(t, err)
    assert.Equal(t, 0, len(errorTodos))
    assert.Contains(t, err.Error(), "database error")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTodos test passed")
}

func TestCreateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTodo")
    fmt.Println("Testing creating a new todo")
    
    // Create a mock todo repository
    mockRepo := new(mocks.MockTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Setup expectations
    mockRepo.On("CreateTodo", "Test Task", "Test Description", true, "user-123").Return(
        models.Todo{
            ID:          "new-todo-id",
            Task:        "Test Task",
            Description: "Test Description",
            Done:        false,
            Important:   true,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        timeStr,
        }, nil)
    
    mockRepo.On("CreateTodo", "Error Task", "Error Description", false, "user-123").Return(
        models.Todo{}, errors.New("failed to create todo"))
    
    // Test successful todo creation
    todo, err := mockRepo.CreateTodo("Test Task", "Test Description", true, "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "new-todo-id", todo.ID)
    assert.Equal(t, "Test Task", todo.Task)
    assert.Equal(t, "Test Description", todo.Description)
    assert.False(t, todo.Done)
    assert.True(t, todo.Important)
    assert.Equal(t, "user-123", todo.UserID)
    assert.Equal(t, dateStr, todo.Date)
    assert.Equal(t, timeStr, todo.Time)
    
    fmt.Printf("Created todo with ID: %s\n", todo.ID)
    fmt.Printf("Todo details: {Task:%s Description:%s Done:%t Important:%t UserID:%s}\n",
        todo.Task, todo.Description, todo.Done, todo.Important, todo.UserID)
    
    // Test error case
    errorTodo, err := mockRepo.CreateTodo("Error Task", "Error Description", false, "user-123")
    
    assert.Error(t, err)
    assert.Empty(t, errorTodo.ID)
    assert.Contains(t, err.Error(), "failed to create todo")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("CreateTodo test passed")
}

func TestUpdateTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateTodo")
    fmt.Println("Testing updating an existing todo")
    
    // Create a mock todo repository
    mockRepo := new(mocks.MockTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Setup expectations
    mockRepo.On("UpdateTodo", "todo-1", "Updated Task", "Updated Description", true, false, "user-123").Return(
        models.Todo{
            ID:          "todo-1",
            Task:        "Updated Task",
            Description: "Updated Description",
            Done:        true,
            Important:   false,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        timeStr,
        }, nil)
    
    mockRepo.On("UpdateTodo", "nonexistent-todo", "Updated Task", "Updated Description", true, false, "user-123").Return(
        models.Todo{}, errors.New("todo not found"))
    
    // Test valid update
    todo, err := mockRepo.UpdateTodo("todo-1", "Updated Task", "Updated Description", true, false, "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "todo-1", todo.ID)
    assert.Equal(t, "Updated Task", todo.Task)
    assert.Equal(t, "Updated Description", todo.Description)
    assert.True(t, todo.Done)
    assert.False(t, todo.Important)
    
    fmt.Printf("Updated todo: {ID:%s Task:%s Description:%s Done:%t Important:%t}\n",
        todo.ID, todo.Task, todo.Description, todo.Done, todo.Important)
    
    // Test update on non-existent todo
    _, err = mockRepo.UpdateTodo("nonexistent-todo", "Updated Task", "Updated Description", true, false, "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent todo: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("UpdateTodo test passed")
}

func TestDeleteTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDeleteTodo")
    fmt.Println("Testing deleting a todo")
    
    // Create a mock todo repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Setup expectations
    mockRepo.On("DeleteTodo", "todo-1", "user-123").Return(nil)
    mockRepo.On("DeleteTodo", "nonexistent-todo", "user-123").Return(errors.New("todo not found"))
    mockRepo.On("DeleteTodo", "todo-2", "wrong-user").Return(errors.New("unauthorized"))
    
    // Test valid deletion
    err := mockRepo.DeleteTodo("todo-1", "user-123")
    assert.NoError(t, err)
    
    fmt.Println("Todo was successfully deleted")
    
    // Test deletion of non-existent todo
    err = mockRepo.DeleteTodo("nonexistent-todo", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent todo: %v\n", err)
    
    // Test deletion with wrong user
    err = mockRepo.DeleteTodo("todo-2", "wrong-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unauthorized")
    
    fmt.Printf("Correctly got error for unauthorized deletion: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("DeleteTodo test passed")
}

func TestUndoTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUndoTodo")
    fmt.Println("Testing undoing a completed todo")
    
    // Create a mock todo repository
    mockRepo := new(mocks.MockTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Setup expectations
    mockRepo.On("UndoTodo", "todo-1", "user-123").Return(
        models.Todo{
            ID:          "todo-1",
            Task:        "Test Task",
            Description: "Test Description",
            Done:        false,
            Important:   true,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        timeStr,
        }, nil)
    
    mockRepo.On("UndoTodo", "nonexistent-todo", "user-123").Return(
        models.Todo{}, errors.New("todo not found"))
    
    // Test undo operation
    todo, err := mockRepo.UndoTodo("todo-1", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "todo-1", todo.ID)
    assert.False(t, todo.Done)
    
    fmt.Printf("Undone todo: {ID:%s Task:%s Done:%t}\n",
        todo.ID, todo.Task, todo.Done)
    
    // Test undo on non-existent todo
    _, err = mockRepo.UndoTodo("nonexistent-todo", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent todo: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("UndoTodo test passed")
}