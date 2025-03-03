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

func TestShareTodoWithUser(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestShareTodoWithUser")
    fmt.Println("Testing sharing a todo with another user")
    
    // Create a mock shared todo repository
    mockRepo := new(mocks.MockSharedTodoRepository)
    
    // Setup expectations
    mockRepo.On("ShareTodoWithUser", "todo-1", "recipient-user", "sender-user").Return(nil)
    mockRepo.On("ShareTodoWithUser", "nonexistent-todo", "recipient-user", "sender-user").Return(
        errors.New("todo not found"))
    mockRepo.On("ShareTodoWithUser", "todo-1", "nonexistent-user", "sender-user").Return(
        errors.New("recipient user not found"))
    mockRepo.On("ShareTodoWithUser", "todo-1", "sender-user", "sender-user").Return(
        errors.New("cannot share todo with yourself"))
    mockRepo.On("ShareTodoWithUser", "already-shared", "recipient-user", "sender-user").Return(
        errors.New("todo is already shared with this user"))
    
    // Test successful sharing
    err := mockRepo.ShareTodoWithUser("todo-1", "recipient-user", "sender-user")
    assert.NoError(t, err)
    
    fmt.Println("Todo successfully shared with recipient user")
    
    // Test sharing non-existent todo
    err = mockRepo.ShareTodoWithUser("nonexistent-todo", "recipient-user", "sender-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "todo not found")
    
    fmt.Printf("Correctly got error for non-existent todo: %v\n", err)
    
    // Test sharing to non-existent user
    err = mockRepo.ShareTodoWithUser("todo-1", "nonexistent-user", "sender-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "recipient user not found")
    
    fmt.Printf("Correctly got error for non-existent recipient: %v\n", err)
    
    // Test sharing with yourself
    err = mockRepo.ShareTodoWithUser("todo-1", "sender-user", "sender-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot share todo with yourself")
    
    fmt.Printf("Correctly got error for sharing with yourself: %v\n", err)
    
    // Test sharing already shared todo
    err = mockRepo.ShareTodoWithUser("already-shared", "recipient-user", "sender-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already shared")
    
    fmt.Printf("Correctly got error for already shared todo: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("ShareTodoWithUser test passed")
}

func TestGetSharedTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetSharedTodos")
    fmt.Println("Testing retrieving todos shared with a user")
    
    // Create a mock shared todo repository
    mockRepo := new(mocks.MockSharedTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Create test shared todos
    sharedTodos := []models.SharedTodo{
        {
            ID:          "shared-todo-1",
            Task:        "Shared Task 1",
            Description: "Shared Description 1",
            Done:        false,
            Important:   true,
            UserID:      "recipient-user",
            Date:        dateStr,
            Time:        timeStr,
            SharedBy:    "sender-user-1",
        },
        {
            ID:          "shared-todo-2",
            Task:        "Shared Task 2",
            Description: "Shared Description 2",
            Done:        true,
            Important:   false,
            UserID:      "recipient-user",
            Date:        dateStr,
            Time:        timeStr,
            SharedBy:    "sender-user-2",
        },
    }
    
    // Setup expectations
    mockRepo.On("GetSharedTodos", "recipient-user").Return(sharedTodos, nil)
    mockRepo.On("GetSharedTodos", "user-no-shared").Return([]models.SharedTodo{}, nil)
    mockRepo.On("GetSharedTodos", "error-user").Return([]models.SharedTodo{}, errors.New("database error"))
    
    // Test getting shared todos for a user with shared todos
    todos, err := mockRepo.GetSharedTodos("recipient-user")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 2, len(todos))
    assert.Equal(t, "Shared Task 1", todos[0].Task)
    assert.Equal(t, "sender-user-1", todos[0].SharedBy)
    assert.Equal(t, "Shared Task 2", todos[1].Task)
    assert.Equal(t, "sender-user-2", todos[1].SharedBy)
    
    fmt.Printf("Retrieved %d shared todos\n", len(todos))
    for i, todo := range todos {
        fmt.Printf("Shared Todo %d: {ID:%s Task:%s SharedBy:%s Done:%t Important:%t}\n",
            i+1, todo.ID, todo.Task, todo.SharedBy, todo.Done, todo.Important)
    }
    
    // Test getting shared todos for a user with no shared todos
    emptyTodos, err := mockRepo.GetSharedTodos("user-no-shared")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyTodos))
    
    fmt.Println("Successfully retrieved empty shared todo list for user with no shared todos")
    
    // Test database error
    _, err = mockRepo.GetSharedTodos("error-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetSharedTodos test passed")
}

func TestGetSharedByMeTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetSharedByMeTodos")
    fmt.Println("Testing retrieving todos shared by a user")
    
    // Create a mock shared todo repository
    mockRepo := new(mocks.MockSharedTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Create test shared todos
    sharedTodos := []models.SharedTodo{
        {
            ID:          "shared-todo-1",
            Task:        "Shared Task 1",
            Description: "Shared Description 1",
            Done:        false,
            Important:   true,
            UserID:      "recipient-user-1",
            Date:        dateStr,
            Time:        timeStr,
            SharedBy:    "sender-user",
        },
        {
            ID:          "shared-todo-2",
            Task:        "Shared Task 2",
            Description: "Shared Description 2",
            Done:        true,
            Important:   false,
            UserID:      "recipient-user-2",
            Date:        dateStr,
            Time:        timeStr,
            SharedBy:    "sender-user",
        },
    }
    
    // Setup expectations
    mockRepo.On("GetSharedByMeTodos", "sender-user").Return(sharedTodos, nil)
    mockRepo.On("GetSharedByMeTodos", "user-no-shared").Return([]models.SharedTodo{}, nil)
    mockRepo.On("GetSharedByMeTodos", "error-user").Return([]models.SharedTodo{}, errors.New("database error"))
    
    // Test getting todos shared by a user who has shared todos
    todos, err := mockRepo.GetSharedByMeTodos("sender-user")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 2, len(todos))
    assert.Equal(t, "Shared Task 1", todos[0].Task)
    assert.Equal(t, "recipient-user-1", todos[0].UserID)
    assert.Equal(t, "Shared Task 2", todos[1].Task)
    assert.Equal(t, "recipient-user-2", todos[1].UserID)
    
    fmt.Printf("Retrieved %d todos shared by user\n", len(todos))
    for i, todo := range todos {
        fmt.Printf("Shared Todo %d: {ID:%s Task:%s SharedWith:%s Done:%t Important:%t}\n",
            i+1, todo.ID, todo.Task, todo.UserID, todo.Done, todo.Important)
    }
    
    // Test getting shared todos for a user who hasn't shared any
    emptyTodos, err := mockRepo.GetSharedByMeTodos("user-no-shared")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyTodos))
    
    fmt.Println("Successfully retrieved empty list for user who hasn't shared any todos")
    
    // Test database error
    _, err = mockRepo.GetSharedByMeTodos("error-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetSharedByMeTodos test passed")
}