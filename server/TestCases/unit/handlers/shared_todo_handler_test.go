package handlers_test

import (
    "context"
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Interfaces for services to allow mocking
type sharedTodoServicer interface {
    GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error)
    GetSharedByMeTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error)
    ShareTodo(ctx context.Context, todoID, recipientID, sharedByID string) error
}
type sharedTodoHandlerTodoServicer interface {
    GetTodoByID(ctx context.Context, id string) (*domain.Todo, error)
}

type sharedTodoUserServicer interface {
    GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}


// getSharedTodosLogic processes the logic for retrieving shared todos
func getSharedTodosLogic(ctx context.Context, userID string, service sharedTodoServicer) (*dto.SharedTodosResponse, error) {
    if userID == "" {
        return nil, errors.New("user ID is required")
    }
    
    // Get todos shared with the user
    received, err := service.GetSharedTodos(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get shared todos: %w", err)
    }
    
    // Get todos shared by the user
    shared, err := service.GetSharedByMeTodos(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get shared-by-me todos: %w", err)
    }
    
    // Combine the responses
    return &dto.SharedTodosResponse{
        Received: received.Received,
        Shared:   shared.Shared,
    }, nil
}

// shareTodoLogic processes the logic for sharing a todo with another user
func shareTodoLogic(ctx context.Context, todoID, recipientUsername, sharedByUserID string, 
    todoService sharedTodoHandlerTodoServicer, userService sharedTodoUserServicer, sharedTodoService sharedTodoServicer) error {
    
    // Validate inputs
    if todoID == "" {
        return errors.New("todo ID is required")
    }
    
    if recipientUsername == "" {
        return errors.New("recipient username is required")
    }
    
    if sharedByUserID == "" {
        return errors.New("shared by user ID is required")
    }
    
    // Get recipient user ID from username
    recipientUser, err := userService.GetUserByUsername(ctx, recipientUsername)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }
    
    // Get todo to verify ownership
    todo, err := todoService.GetTodoByID(ctx, todoID)
    if err != nil {
        return fmt.Errorf("todo not found: %w", err)
    }
    
    // Verify ownership
    if todo.UserID != sharedByUserID {
        return errors.New("unauthorized: you can only share your own todos")
    }
    
    // Share the todo
    return sharedTodoService.ShareTodo(ctx, todoID, recipientUser.ID, sharedByUserID)
}

func TestGetSharedTodosLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetSharedTodosLogic ===")
    fmt.Println("Testing get shared todos logic")
    
    // Create mock service
    mockService := new(mocks.MockSharedTodoService)
    
    // Setup test data
    userID := "user-123"
    currentTime := time.Now()
    
    // Create sample received todos
    receivedTodos := []dto.SharedTodoResponse{
        {
            ID:          "shared-todo-1",
            Task:        "Shared Task 1",
            Description: "Description 1",
            Done:        false,
            Important:   true,
            UserID:      userID,
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    "sender-user-1",
        },
        {
            ID:          "shared-todo-2",
            Task:        "Shared Task 2",
            Description: "Description 2",
            Done:        true,
            Important:   false,
            UserID:      userID,
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    "sender-user-2",
        },
    }

    // Create sample shared by me todos
    sharedByMeTodos := []dto.SharedTodoResponse{
        {
            ID:          "shared-by-me-1",
            Task:        "My Shared Task 1",
            Description: "My Description 1",
            Done:        false,
            Important:   true,
            UserID:      "recipient-user-1",
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    userID,
        },
        {
            ID:          "shared-by-me-2",
            Task:        "My Shared Task 2",
            Description: "My Description 2",
            Done:        true,
            Important:   false,
            UserID:      "recipient-user-2",
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    userID,
        },
    }
    
    // Setup expectations for successful scenario
    mockService.On("GetSharedTodos", mock.Anything, userID).Return(&dto.SharedTodosResponse{Received: receivedTodos}, nil)
    mockService.On("GetSharedByMeTodos", mock.Anything, userID).Return(&dto.SharedTodosResponse{Shared: sharedByMeTodos}, nil)
    
    // Scenario 1: Test successful retrieval
    fmt.Println("Scenario 1: Testing successful shared todos retrieval")
    response, err := getSharedTodosLogic(context.Background(), userID, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, 2, len(response.Received))
    assert.Equal(t, 2, len(response.Shared))
    assert.Equal(t, "Shared Task 1", response.Received[0].Task)
    assert.Equal(t, "My Shared Task 1", response.Shared[0].Task)
    assert.Equal(t, "sender-user-1", response.Received[0].SharedBy)
    assert.Equal(t, "recipient-user-1", response.Shared[0].UserID)
    fmt.Printf("✅ Successfully retrieved %d received todos and %d shared todos\n", 
        len(response.Received), len(response.Shared))
    
    // Setup expectations for error scenario - GetSharedTodos fails
    mockService.On("GetSharedTodos", mock.Anything, "error-user").Return(nil, errors.New("database error"))
    
    // Scenario 2: Test error when retrieving shared todos
    fmt.Println("\nScenario 2: Testing database error when retrieving shared todos")
    response, err = getSharedTodosLogic(context.Background(), "error-user", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, response)
    assert.Contains(t, err.Error(), "failed to get shared todos")
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations for another error scenario - GetSharedByMeTodos fails
    mockService.On("GetSharedTodos", mock.Anything, "error-user-2").Return(&dto.SharedTodosResponse{Received: receivedTodos}, nil)
    mockService.On("GetSharedByMeTodos", mock.Anything, "error-user-2").Return(nil, errors.New("database error"))
    
    // Scenario 3: Test error when retrieving shared-by-me todos
    fmt.Println("\nScenario 3: Testing database error when retrieving shared by me todos")
    response, err = getSharedTodosLogic(context.Background(), "error-user-2", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, response)
    assert.Contains(t, err.Error(), "failed to get shared-by-me todos")
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 4: Test empty user ID
    fmt.Println("\nScenario 4: Testing empty user ID")
    response, err = getSharedTodosLogic(context.Background(), "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, response)
    assert.Contains(t, err.Error(), "user ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All GetSharedTodosLogic test scenarios passed")
}

func TestShareTodoLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestShareTodoLogic ===")
    fmt.Println("Testing sharing a todo logic")
    
    // Create mock services
    mockSharedTodoService := new(mocks.MockSharedTodoService)
    mockUserService := new(mocks.MockUserService)
    mockTodoService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    recipientUsername := "recipient-user"
    recipientID := "recipient-123"
    
    // Setup user to be returned by user service
    mockRecipientUser := &domain.User{
        ID:       recipientID,
        Username: recipientUsername,
        Password: "hashed_password",
    }
    
    // Setup todo to be returned by todo service
    mockTodo := &domain.Todo{
        ID:          todoID,
        Task:        "Test Task",
        Description: "Test Description",
        Done:        false,
        Important:   true,
        UserID:      userID, // Owner is the current user
        Date:        time.Now(),
        Time:        time.Now(),
    }
    
    // Setup expectations
    mockUserService.On("GetUserByUsername", mock.Anything, recipientUsername).Return(mockRecipientUser, nil)
    mockTodoService.On("GetTodoByID", mock.Anything, todoID).Return(mockTodo, nil)
    mockSharedTodoService.On("ShareTodo", mock.Anything, todoID, recipientID, userID).Return(nil)
    
    // Scenario 1: Test successful sharing
    fmt.Println("Scenario 1: Testing successful todo sharing")
    err := shareTodoLogic(context.Background(), todoID, recipientUsername, userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.NoError(t, err)
    fmt.Printf("✅ Todo %s successfully shared with user %s\n", todoID, recipientUsername)
    
    // Setup expectations for non-existent user
    mockUserService.On("GetUserByUsername", mock.Anything, "nonexistent-user").Return(nil, errors.New("user not found"))
    
    // Scenario 2: Test non-existent recipient
    fmt.Println("\nScenario 2: Testing sharing with non-existent user")
    err = shareTodoLogic(context.Background(), todoID, "nonexistent-user", userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "user not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations for non-existent todo
    mockTodoService.On("GetTodoByID", mock.Anything, "nonexistent-todo").Return(nil, errors.New("todo not found"))
    
    // Scenario 3: Test non-existent todo
    fmt.Println("\nScenario 3: Testing sharing non-existent todo")
    err = shareTodoLogic(context.Background(), "nonexistent-todo", recipientUsername, userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "todo not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup todo not owned by user
    notOwnedTodo := &domain.Todo{
        ID:          "not-owned-todo",
        Task:        "Not My Task",
        Description: "Not My Description",
        Done:        false,
        Important:   true,
        UserID:      "different-user", // Not the current user
        Date:        time.Now(),
        Time:        time.Now(),
    }
    mockTodoService.On("GetTodoByID", mock.Anything, "not-owned-todo").Return(notOwnedTodo, nil)
    
    // Scenario 4: Test not owner of todo
    fmt.Println("\nScenario 4: Testing sharing a todo not owned by the user")
    err = shareTodoLogic(context.Background(), "not-owned-todo", recipientUsername, userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup already shared todo
    alreadySharedTodo := &domain.Todo{
        ID:          "already-shared-todo",
        Task:        "Already Shared Task",
        Description: "Already Shared Description",
        Done:        false,
        Important:   true,
        UserID:      userID, // Owner is the current user
        Date:        time.Now(),
        Time:        time.Now(),
    }
    mockTodoService.On("GetTodoByID", mock.Anything, "already-shared-todo").Return(alreadySharedTodo, nil)
    mockSharedTodoService.On("ShareTodo", mock.Anything, "already-shared-todo", recipientID, userID).Return(errors.New("todo already shared with this user"))
    
    // Scenario 5: Test already shared todo
    fmt.Println("\nScenario 5: Testing sharing a todo that is already shared")
    err = shareTodoLogic(context.Background(), "already-shared-todo", recipientUsername, userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already shared")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 6: Test empty todo ID
    fmt.Println("\nScenario 6: Testing sharing with empty todo ID")
    err = shareTodoLogic(context.Background(), "", recipientUsername, userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "todo ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 7: Test empty recipient username
    fmt.Println("\nScenario 7: Testing sharing with empty recipient username")
    err = shareTodoLogic(context.Background(), todoID, "", userID, mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "recipient username is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 8: Test empty user ID
    fmt.Println("\nScenario 8: Testing sharing with empty user ID")
    err = shareTodoLogic(context.Background(), todoID, recipientUsername, "", mockTodoService, mockUserService, mockSharedTodoService)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "shared by user ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockSharedTodoService.AssertExpectations(t)
    mockUserService.AssertExpectations(t)
    mockTodoService.AssertExpectations(t)
    fmt.Println("✅ All ShareTodoLogic test scenarios passed")
}