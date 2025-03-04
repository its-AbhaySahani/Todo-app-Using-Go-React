package services_test

import (
    "context"
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/shared_todos"
    "github.com/stretchr/testify/assert"
)

func TestGetSharedTodos(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetSharedTodos ===")
    fmt.Println("Testing retrieving todos shared with a user")
    
    // Create mock repositories
    mockRepo := new(mocks.MockSharedTodoRepository)
    mockTodoRepo := new(mocks.MockTodoRepository)
    mockUserRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repositories
    service := shared_todos.NewSharedTodoService(mockRepo, mockTodoRepo, mockUserRepo)
    
    // Setup test data
    userID := "user-123"
    currentTime := time.Now()
    
    // Create test shared todos
    sharedTodos := []domain.SharedTodo{
        {
            ID:          "shared-todo-1",
            Task:        "Shared Task 1",
            Description: "Shared Description 1",
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
            Description: "Shared Description 2",
            Done:        true,
            Important:   false,
            UserID:      userID,
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    "sender-user-2",
        },
    }
    
    // Set up expectations
    mockRepo.On("GetSharedTodos", context.Background(), userID).Return(sharedTodos, nil)
    mockRepo.On("GetSharedTodos", context.Background(), "empty-user").Return([]domain.SharedTodo{}, nil)
    mockRepo.On("GetSharedTodos", context.Background(), "error-user").Return([]domain.SharedTodo{}, errors.New("database error"))
    
    // Scenario 1: User with shared todos
    fmt.Println("Scenario 1: Testing retrieving shared todos for a user")
    res, err := service.GetSharedTodos(context.Background(), userID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 2, len(res.Received))
    fmt.Printf("✅ Successfully retrieved %d shared todos for user %s\n", len(res.Received), userID)
    fmt.Println("   Shared todo details:")
    for i, todo := range res.Received {
        fmt.Printf("   - Todo %d: {ID: %s, Task: %s, SharedBy: %s}\n", 
            i+1, todo.ID, todo.Task, todo.SharedBy)
    }
    
    // Scenario 2: User with no shared todos
    fmt.Println("\nScenario 2: Testing retrieving shared todos for a user with no shared todos")
    res, err = service.GetSharedTodos(context.Background(), "empty-user")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 0, len(res.Received))
    fmt.Println("✅ Successfully returned empty list for user with no shared todos")
    
    // Scenario 3: Database error
    fmt.Println("\nScenario 3: Testing database error when retrieving shared todos")
    res, err = service.GetSharedTodos(context.Background(), "error-user")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All GetSharedTodos test scenarios passed")
}

func TestShareTodo(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestShareTodo ===")
    fmt.Println("Testing sharing a todo with another user")
    
    // Create mock repositories
    mockRepo := new(mocks.MockSharedTodoRepository)
    mockTodoRepo := new(mocks.MockTodoRepository)
    mockUserRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repositories
    service := shared_todos.NewSharedTodoService(mockRepo, mockTodoRepo, mockUserRepo)
    
    // Setup test data
    todoID := "todo-123"
    ownerID := "owner-123"
    recipientID := "recipient-123"
    currentTime := time.Now()
    
    // Mock the original todo
    originalTodo := domain.Todo{
        ID:          todoID,
        Task:        "Original Task",
        Description: "Original Description",
        Done:        false,
        Important:   true,
        UserID:      ownerID,
        Date:        currentTime,
        Time:        currentTime,
    }
    
    // Also create a todo for the "already shared" case
    alreadySharedTodo := domain.Todo{
        ID:          "already-shared",
        Task:        "Already Shared Task",
        Description: "Already Shared Description",
        Done:        false,
        Important:   true,
        UserID:      ownerID,  // Same owner as original todo
        Date:        currentTime,
        Time:        currentTime,
    }
    
    // Set up expectations
    mockTodoRepo.On("GetTodoByID", context.Background(), todoID).Return(&originalTodo, nil)
    mockTodoRepo.On("GetTodoByID", context.Background(), "nonexistent-todo").Return(nil, errors.New("todo not found"))
    mockTodoRepo.On("GetTodoByID", context.Background(), "not-owned-todo").Return(&domain.Todo{
        ID:     "not-owned-todo",
        UserID: "different-owner",
    }, nil)
    // Add expectation for the already shared todo
    mockTodoRepo.On("GetTodoByID", context.Background(), "already-shared").Return(&alreadySharedTodo, nil)
    
    mockRepo.On("IsSharedWithUser", context.Background(), todoID, recipientID).Return(false, nil)
    mockRepo.On("IsSharedWithUser", context.Background(), "already-shared", recipientID).Return(true, nil)
    
    mockRepo.On("ShareTodo", context.Background(), todoID, recipientID, ownerID).Return(nil)
    
    // Scenario 1: Successful share
    fmt.Println("Scenario 1: Testing successful todo sharing")
    err := service.ShareTodo(context.Background(), todoID, recipientID, ownerID)
    
    // Assertions
    assert.NoError(t, err)
    fmt.Printf("✅ Todo %s successfully shared with user %s by user %s\n", todoID, recipientID, ownerID)
    
    // Scenario 2: Todo not found
    fmt.Println("\nScenario 2: Testing sharing non-existent todo")
    err = service.ShareTodo(context.Background(), "nonexistent-todo", recipientID, ownerID)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to get todo")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 3: Not owner of todo
    fmt.Println("\nScenario 3: Testing sharing a todo that user doesn't own")
    err = service.ShareTodo(context.Background(), "not-owned-todo", recipientID, ownerID)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 4: Already shared todo
    fmt.Println("\nScenario 4: Testing sharing a todo that is already shared")
    err = service.ShareTodo(context.Background(), "already-shared", recipientID, ownerID)
    
    // Assertions
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already shared")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    mockTodoRepo.AssertExpectations(t)
    
    fmt.Println("✅ All ShareTodo test scenarios passed")
}

func TestGetSharedByMeTodos(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetSharedByMeTodos ===")
    fmt.Println("Testing retrieving todos shared by a user")
    
    // Create mock repositories
    mockRepo := new(mocks.MockSharedTodoRepository)
    mockTodoRepo := new(mocks.MockTodoRepository)
    mockUserRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repositories
    service := shared_todos.NewSharedTodoService(mockRepo, mockTodoRepo, mockUserRepo)
    
    // Setup test data
    userID := "user-123"
    currentTime := time.Now()
    
    // Create test shared todos
    sharedByMeTodos := []domain.SharedTodo{
        {
            ID:          "shared-todo-1",
            Task:        "Shared Task 1",
            Description: "Shared Description 1",
            Done:        false,
            Important:   true,
            UserID:      "recipient-user-1",
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    userID,
        },
        {
            ID:          "shared-todo-2",
            Task:        "Shared Task 2",
            Description: "Shared Description 2",
            Done:        true,
            Important:   false,
            UserID:      "recipient-user-2",
            Date:        currentTime,
            Time:        currentTime,
            SharedBy:    userID,
        },
    }
    
    // Set up expectations
    mockRepo.On("GetSharedByMeTodos", context.Background(), userID).Return(sharedByMeTodos, nil)
    mockRepo.On("GetSharedByMeTodos", context.Background(), "empty-user").Return([]domain.SharedTodo{}, nil)
    mockRepo.On("GetSharedByMeTodos", context.Background(), "error-user").Return([]domain.SharedTodo{}, errors.New("database error"))
    
    // Scenario 1: User who has shared todos with others
    fmt.Println("Scenario 1: Testing retrieving todos shared by a user")
    res, err := service.GetSharedByMeTodos(context.Background(), userID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 2, len(res.Shared))
    fmt.Printf("✅ Successfully retrieved %d todos shared by user %s\n", len(res.Shared), userID)
    fmt.Println("   Shared todo details:")
    for i, todo := range res.Shared {
        fmt.Printf("   - Todo %d: {ID: %s, Task: %s, SharedWith: %s}\n", 
            i+1, todo.ID, todo.Task, todo.UserID)
    }
    
    // Scenario 2: User who hasn't shared any todos
    fmt.Println("\nScenario 2: Testing retrieving shared todos for a user who hasn't shared any")
    res, err = service.GetSharedByMeTodos(context.Background(), "empty-user")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 0, len(res.Shared))
    fmt.Println("✅ Successfully returned empty list for user who hasn't shared any todos")
    
    // Scenario 3: Database error
    fmt.Println("\nScenario 3: Testing database error when retrieving shared todos")
    res, err = service.GetSharedByMeTodos(context.Background(), "error-user")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All GetSharedByMeTodos test scenarios passed")
}