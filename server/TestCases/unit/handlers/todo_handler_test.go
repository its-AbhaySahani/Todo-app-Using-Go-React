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

// Interface for todo service to allow mocking
type todoServicer interface {
    CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error)
    GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error)
    UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error)
    DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error)
    UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error)
    GetTodoByID(ctx context.Context, id string) (*domain.Todo, error)
}

// getTodosLogic processes the logic for retrieving a user's todos
func getTodosLogic(ctx context.Context, userID string, service todoServicer) (*dto.TodosResponse, error) {
    if userID == "" {
        return nil, errors.New("user ID is required")
    }
    
    return service.GetTodosByUserID(ctx, userID)
}

// createTodoLogic processes the logic for creating a new todo
func createTodoLogic(ctx context.Context, req *dto.CreateTodoRequest, service todoServicer) (*dto.CreateResponse, error) {
    if req.Task == "" {
        return nil, errors.New("task is required")
    }
    
    return service.CreateTodo(ctx, req)
}

// updateTodoLogic processes the logic for updating a todo
func updateTodoLogic(ctx context.Context, req *dto.UpdateTodoRequest, service todoServicer) (*dto.SuccessResponse, error) {
    if req.ID == "" {
        return nil, errors.New("todo ID is required")
    }
    
    if req.Task == "" {
        return nil, errors.New("task is required")
    }
    
    return service.UpdateTodo(ctx, req)
}

// deleteTodoLogic processes the logic for deleting a todo
func deleteTodoLogic(ctx context.Context, todoID, userID string, service todoServicer) (*dto.SuccessResponse, error) {
    if todoID == "" {
        return nil, errors.New("todo ID is required")
    }
    
    if userID == "" {
        return nil, errors.New("user ID is required")
    }
    
    return service.DeleteTodo(ctx, todoID, userID)
}

// undoTodoLogic processes the logic for undoing a todo
func undoTodoLogic(ctx context.Context, todoID, userID string, service todoServicer) (*dto.SuccessResponse, error) {
    if todoID == "" {
        return nil, errors.New("todo ID is required")
    }
    
    if userID == "" {
        return nil, errors.New("user ID is required")
    }
    
    return service.UndoTodo(ctx, todoID, userID)
}

func TestGetTodosLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetTodosLogic ===")
    fmt.Println("Testing get todos logic")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    today := time.Now()
    
    // Create sample todos
    todos := []dto.TodoResponse{
        {
            ID:          "todo-1",
            Task:        "Test Task 1",
            Description: "Description 1",
            Done:        false,
            Important:   true,
            UserID:      userID,
            Date:        today,
            Time:        today,
        },
        {
            ID:          "todo-2",
            Task:        "Test Task 2",
            Description: "Description 2",
            Done:        true,
            Important:   false,
            UserID:      userID,
            Date:        today,
            Time:        today,
        },
    }
    
    // Setup expectations for success
    mockService.On("GetTodosByUserID", mock.Anything, userID).Return(&dto.TodosResponse{Todos: todos}, nil)
    
    // Scenario 1: Test successful retrieval
    fmt.Println("Scenario 1: Testing successful todos retrieval")
    res, err := getTodosLogic(context.Background(), userID, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, 2, len(res.Todos))
    assert.Equal(t, "todo-1", res.Todos[0].ID)
    assert.Equal(t, "Test Task 1", res.Todos[0].Task)
    fmt.Printf("✅ Successfully retrieved %d todos for user\n", len(res.Todos))
    
    // Setup expectations for database error
    mockService.On("GetTodosByUserID", mock.Anything, "error-user").Return(nil, errors.New("database error"))
    
    // Scenario 2: Test database error
    fmt.Println("\nScenario 2: Testing database error when retrieving todos")
    res, err = getTodosLogic(context.Background(), "error-user", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 3: Test empty user ID
    fmt.Println("\nScenario 3: Testing empty user ID")
    res, err = getTodosLogic(context.Background(), "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "user ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All GetTodosLogic test scenarios passed")
}

func TestCreateTodoLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestCreateTodoLogic ===")
    fmt.Println("Testing create todo logic")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    task := "New Task"
    description := "New Description"
    important := true
    todoID := "new-todo-1"
    date := time.Now()
    todoTime := time.Now()
    
    // Create request
    req := &dto.CreateTodoRequest{
        Task:        task,
        Description: description,
        Important:   important,
        UserID:      userID,
        Date:        date,
        Time:        todoTime,
    }
    
    // Setup expectations for success
    mockService.On(
        "CreateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(r *dto.CreateTodoRequest) bool {
            return r.Task == task && 
                   r.Description == description && 
                   r.Important == important &&
                   r.UserID == userID
        }),
    ).Return(&dto.CreateResponse{ID: todoID}, nil)
    
    // Scenario 1: Test successful creation
    fmt.Println("Scenario 1: Testing successful todo creation")
    res, err := createTodoLogic(context.Background(), req, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, todoID, res.ID)
    fmt.Printf("✅ Todo created successfully with ID: %s\n", res.ID)
    
    // Setup expectations for database error
    mockService.On(
        "CreateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(r *dto.CreateTodoRequest) bool {
            return r.Task == "Error Task"
        }),
    ).Return(nil, errors.New("database error"))
    
    // Create error request
    errorReq := &dto.CreateTodoRequest{
        Task:        "Error Task",
        Description: "Error Description",
        Important:   false,
        UserID:      userID,
    }
    
    // Scenario 2: Test database error
    fmt.Println("\nScenario 2: Testing creation with database error")
    res, err = createTodoLogic(context.Background(), errorReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "database error")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty task
    emptyTaskReq := &dto.CreateTodoRequest{
        Task:        "",
        Description: "Some description",
        Important:   false,
        UserID:      userID,
    }
    
    // Scenario 3: Test empty task
    fmt.Println("\nScenario 3: Testing creation with empty task")
    res, err = createTodoLogic(context.Background(), emptyTaskReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "task is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All CreateTodoLogic test scenarios passed")
}

func TestUpdateTodoLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUpdateTodoLogic ===")
    fmt.Println("Testing update todo logic")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    task := "Updated Task"
    description := "Updated Description"
    done := true
    important := false
    
    // Create request
    req := &dto.UpdateTodoRequest{
        ID:          todoID,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    // Setup expectations for success
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(r *dto.UpdateTodoRequest) bool {
            return r.ID == todoID && 
                   r.Task == task && 
                   r.Description == description && 
                   r.Done == done &&
                   r.Important == important &&
                   r.UserID == userID
        }),
    ).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Scenario 1: Test successful update
    fmt.Println("Scenario 1: Testing successful todo update")
    res, err := updateTodoLogic(context.Background(), req, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s updated successfully\n", todoID)
    
    // Setup expectations for non-existent todo
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(r *dto.UpdateTodoRequest) bool {
            return r.ID == "nonexistent-todo"
        }),
    ).Return(nil, errors.New("todo not found"))
    
    // Create non-existent todo request
    nonExistentReq := &dto.UpdateTodoRequest{
        ID:          "nonexistent-todo",
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    // Scenario 2: Test non-existent todo
    fmt.Println("\nScenario 2: Testing update of non-existent todo")
    res, err = updateTodoLogic(context.Background(), nonExistentReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations for unauthorized update
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(r *dto.UpdateTodoRequest) bool {
            return r.ID == todoID && r.UserID == "wrong-user"
        }),
    ).Return(nil, errors.New("unauthorized"))
    
    // Create unauthorized request
    unauthorizedReq := &dto.UpdateTodoRequest{
        ID:          todoID,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      "wrong-user",
    }
    
    // Scenario 3: Test unauthorized update
    fmt.Println("\nScenario 3: Testing unauthorized update")
    res, err = updateTodoLogic(context.Background(), unauthorizedReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty ID
    emptyIDReq := &dto.UpdateTodoRequest{
        ID:          "",
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    // Scenario 4: Test empty ID
    fmt.Println("\nScenario 4: Testing update with empty ID")
    res, err = updateTodoLogic(context.Background(), emptyIDReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty task
    emptyTaskReq := &dto.UpdateTodoRequest{
        ID:          todoID,
        Task:        "",
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    // Scenario 5: Test empty task
    fmt.Println("\nScenario 5: Testing update with empty task")
    res, err = updateTodoLogic(context.Background(), emptyTaskReq, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "task is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All UpdateTodoLogic test scenarios passed")
}

func TestDeleteTodoLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestDeleteTodoLogic ===")
    fmt.Println("Testing delete todo logic")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    
    // Setup expectations for success
    mockService.On("DeleteTodo", mock.Anything, todoID, userID).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Scenario 1: Test successful deletion
    fmt.Println("Scenario 1: Testing successful todo deletion")
    res, err := deleteTodoLogic(context.Background(), todoID, userID, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s deleted successfully\n", todoID)
    
    // Setup expectations for non-existent todo
    mockService.On("DeleteTodo", mock.Anything, "nonexistent-todo", userID).Return(nil, errors.New("todo not found"))
    
    // Scenario 2: Test non-existent todo
    fmt.Println("\nScenario 2: Testing deletion of non-existent todo")
    res, err = deleteTodoLogic(context.Background(), "nonexistent-todo", userID, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations for unauthorized deletion
    mockService.On("DeleteTodo", mock.Anything, todoID, "wrong-user").Return(nil, errors.New("unauthorized"))
    
    // Scenario 3: Test unauthorized deletion
    fmt.Println("\nScenario 3: Testing unauthorized deletion")
    res, err = deleteTodoLogic(context.Background(), todoID, "wrong-user", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 4: Test empty todo ID
    fmt.Println("\nScenario 4: Testing deletion with empty todo ID")
    res, err = deleteTodoLogic(context.Background(), "", userID, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 5: Test empty user ID
    fmt.Println("\nScenario 5: Testing deletion with empty user ID")
    res, err = deleteTodoLogic(context.Background(), todoID, "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "user ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All DeleteTodoLogic test scenarios passed")
}

func TestUndoTodoLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUndoTodoLogic ===")
    fmt.Println("Testing undo todo logic")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    
    // Setup expectations for success
    mockService.On("UndoTodo", mock.Anything, todoID, userID).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Scenario 1: Test successful undoing
    fmt.Println("Scenario 1: Testing successful todo undoing")
    res, err := undoTodoLogic(context.Background(), todoID, userID, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.True(t, res.Success)
    fmt.Printf("✅ Todo %s undone successfully\n", todoID)
    
    // Setup expectations for non-existent todo
    mockService.On("UndoTodo", mock.Anything, "nonexistent-todo", userID).Return(nil, errors.New("todo not found"))
    
    // Scenario 2: Test non-existent todo
    fmt.Println("\nScenario 2: Testing undoing of non-existent todo")
    res, err = undoTodoLogic(context.Background(), "nonexistent-todo", userID, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations for unauthorized undoing
    mockService.On("UndoTodo", mock.Anything, todoID, "wrong-user").Return(nil, errors.New("unauthorized"))
    
    // Scenario 3: Test unauthorized undoing
    fmt.Println("\nScenario 3: Testing unauthorized undoing")
    res, err = undoTodoLogic(context.Background(), todoID, "wrong-user", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "unauthorized")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 4: Test empty todo ID
    fmt.Println("\nScenario 4: Testing undoing with empty todo ID")
    res, err = undoTodoLogic(context.Background(), "", userID, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "todo ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Scenario 5: Test empty user ID
    fmt.Println("\nScenario 5: Testing undoing with empty user ID")
    res, err = undoTodoLogic(context.Background(), todoID, "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "user ID is required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All UndoTodoLogic test scenarios passed")
}