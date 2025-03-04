package services_test

import (
    "context"
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/todos"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreateTodo(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestCreateTodo ===")
    fmt.Println("Testing todo creation functionality")
    
    // Create a mock repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create the service with the mock repository
    todoService := todos.NewTodoService(mockRepo)
    
    // Setup test data
    todoID := "todo-123"
    taskName := "Test Task"
    description := "Test Description"
    userID := "user-123"
    
    // Mock expected date and time
    date := time.Now()
    todoTime := time.Date(2000, 1, 1, 10, 30, 0, 0, time.UTC) // Use valid year (2000)
    
    // Set up expectations on the mock repository
    mockRepo.On("CreateTodo", 
        context.Background(), 
        taskName, 
        description, 
        false, // done
        true,  // important
        userID,
        mock.AnythingOfType("time.Time"),
        mock.AnythingOfType("time.Time"),
    ).Return(todoID, nil)
    
    // Create the request
    req := &dto.CreateTodoRequest{
        Task:        taskName,
        Description: description,
        Done:        false,
        Important:   true,
        UserID:      userID,
        Date:        date,
        Time:        todoTime,
    }
    
    // Scenario 1: Successful creation
    fmt.Println("Scenario 1: Testing successful todo creation")
    res, err := todoService.CreateTodo(context.Background(), req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, todoID, res.ID)
    fmt.Printf("✅ Todo created successfully with ID: %s\n", res.ID)
    
    // Scenario 2: Database error
    fmt.Println("\nScenario 2: Testing todo creation with database error")
    mockRepo.On("CreateTodo", 
        context.Background(), 
        "Error Task", 
        "Error Description", 
        false,
        false,
        userID,
        mock.AnythingOfType("time.Time"),
        mock.AnythingOfType("time.Time"),
    ).Return("", errors.New("database error"))
    
    // Create a request that will cause an error
    errorReq := &dto.CreateTodoRequest{
        Task:        "Error Task",
        Description: "Error Description",
        Done:        false,
        Important:   false,
        UserID:      userID,
        Date:        date,
        Time:        todoTime,
    }
    
    // Test the service with error case
    res, err = todoService.CreateTodo(context.Background(), errorReq)
    
    // Assertions for error case
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All CreateTodo test scenarios passed")
}

func TestGetTodosByUserID(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetTodosByUserID ===")
    fmt.Println("Testing retrieving todos for a user")
    
    // Create a mock repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create the service with the mock repository
    todoService := todos.NewTodoService(mockRepo)
    
    // Setup test data
    userID := "user-123"
    mockTodos := []domain.Todo{
        {
            ID:          "todo-1",
            Task:        "Task 1",
            Description: "Description 1",
            Done:        false,
            Important:   true,
            UserID:      userID,
            Date:        time.Now(),
            Time:        time.Now(),
        },
        {
            ID:          "todo-2",
            Task:        "Task 2",
            Description: "Description 2",
            Done:        true,
            Important:   false,
            UserID:      userID,
            Date:        time.Now(),
            Time:        time.Now(),
        },
    }
    
    // Set up expectations
    mockRepo.On("GetTodosByUserID", context.Background(), userID).Return(mockTodos, nil)
    mockRepo.On("GetTodosByUserID", context.Background(), "empty-user").Return([]domain.Todo{}, nil)
    mockRepo.On("GetTodosByUserID", context.Background(), "error-user").Return([]domain.Todo{}, errors.New("database error"))
    
    // Scenario 1: User with todos
    fmt.Println("Scenario 1: Testing retrieving todos for a user with todos")
    res, err := todoService.GetTodosByUserID(context.Background(), userID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 2, len(res.Todos))
    fmt.Printf("✅ Successfully retrieved %d todos for user %s\n", len(res.Todos), userID)
    fmt.Println("   Todo details:")
    for i, todo := range res.Todos {
        fmt.Printf("   - Todo %d: {ID: %s, Task: %s, Done: %t, Important: %t}\n", 
            i+1, todo.ID, todo.Task, todo.Done, todo.Important)
    }
    
    // Scenario 2: User with no todos
    fmt.Println("\nScenario 2: Testing retrieving todos for a user with no todos")
    res, err = todoService.GetTodosByUserID(context.Background(), "empty-user")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 0, len(res.Todos))
    fmt.Println("✅ Successfully returned empty list for user with no todos")
    
    // Scenario 3: Database error
    fmt.Println("\nScenario 3: Testing database error when retrieving todos")
    res, err = todoService.GetTodosByUserID(context.Background(), "error-user")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All GetTodosByUserID test scenarios passed")
}

func TestUpdateTodo(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUpdateTodo ===")
    fmt.Println("Testing todo update functionality")
    
    // Create a mock repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create the service with the mock repository
    todoService := todos.NewTodoService(mockRepo)
    
    // Setup test data
    todoID := "todo-123"
    userID := "user-123"
    task := "Updated Task"
    description := "Updated Description"
    done := true
    important := false
    
    // Set up expectations
    mockRepo.On("UpdateTodo", 
        context.Background(),
        todoID,
        task,
        description,
        done,
        important,
        userID,
    ).Return(true, nil)
    
    mockRepo.On("UpdateTodo", 
        context.Background(),
        "nonexistent-todo",
        task,
        description,
        done,
        important,
        userID,
    ).Return(false, errors.New("todo not found"))
    
    mockRepo.On("UpdateTodo", 
        context.Background(),
        todoID,
        task,
        description,
        done,
        important,
        "wrong-user",
    ).Return(false, errors.New("unauthorized"))
    
    // Create the request
    req := &dto.UpdateTodoRequest{
        ID:          todoID,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    // Scenario 1: Successful update
    fmt.Println("Scenario 1: Testing successful todo update")
    res, err := todoService.UpdateTodo(context.Background(), req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s updated successfully\n", todoID)
    
    // Scenario 2: Non-existent todo
    fmt.Println("\nScenario 2: Testing update of non-existent todo")
    req.ID = "nonexistent-todo"
    res, err = todoService.UpdateTodo(context.Background(), req)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 3: Unauthorized update
    fmt.Println("\nScenario 3: Testing unauthorized update")
    req.ID = todoID
    req.UserID = "wrong-user"
    res, err = todoService.UpdateTodo(context.Background(), req)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All UpdateTodo test scenarios passed")
}

func TestDeleteTodo(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestDeleteTodo ===")
    fmt.Println("Testing todo deletion functionality")
    
    // Create a mock repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create the service with the mock repository
    todoService := todos.NewTodoService(mockRepo)
    
    // Setup test data
    todoID := "todo-123"
    userID := "user-123"
    
    // Set up expectations
    mockRepo.On("DeleteTodo", context.Background(), todoID, userID).Return(true, nil)
    mockRepo.On("DeleteTodo", context.Background(), "nonexistent-todo", userID).Return(false, errors.New("todo not found"))
    mockRepo.On("DeleteTodo", context.Background(), todoID, "wrong-user").Return(false, errors.New("unauthorized"))
    
    // Scenario 1: Successful deletion
    fmt.Println("Scenario 1: Testing successful todo deletion")
    res, err := todoService.DeleteTodo(context.Background(), todoID, userID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s deleted successfully\n", todoID)
    
    // Scenario 2: Non-existent todo
    fmt.Println("\nScenario 2: Testing deletion of non-existent todo")
    res, err = todoService.DeleteTodo(context.Background(), "nonexistent-todo", userID)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 3: Unauthorized deletion
    fmt.Println("\nScenario 3: Testing unauthorized deletion")
    res, err = todoService.DeleteTodo(context.Background(), todoID, "wrong-user")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All DeleteTodo test scenarios passed")
}

func TestUndoTodo(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUndoTodo ===")
    fmt.Println("Testing todo undo functionality")
    
    // Create a mock repository
    mockRepo := new(mocks.MockTodoRepository)
    
    // Create the service with the mock repository
    todoService := todos.NewTodoService(mockRepo)
    
    // Setup test data
    todoID := "todo-123"
    userID := "user-123"
    
    // Set up expectations
    mockRepo.On("UndoTodo", context.Background(), todoID, userID).Return(true, nil)
    mockRepo.On("UndoTodo", context.Background(), "nonexistent-todo", userID).Return(false, errors.New("todo not found"))
    mockRepo.On("UndoTodo", context.Background(), todoID, "wrong-user").Return(false, errors.New("unauthorized"))
    
    // Scenario 1: Successful undo
    fmt.Println("Scenario 1: Testing successful todo undo")
    res, err := todoService.UndoTodo(context.Background(), todoID, userID)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s undone successfully\n", todoID)
    
    // Scenario 2: Non-existent todo
    fmt.Println("\nScenario 2: Testing undo of non-existent todo")
    res, err = todoService.UndoTodo(context.Background(), "nonexistent-todo", userID)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 3: Unauthorized undo
    fmt.Println("\nScenario 3: Testing unauthorized undo")
    res, err = todoService.UndoTodo(context.Background(), todoID, "wrong-user")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All UndoTodo test scenarios passed")
}