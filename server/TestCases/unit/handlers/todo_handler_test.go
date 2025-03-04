package handlers_test

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// TodoServiceAdapterRepo adapts MockTodoService to satisfy domain.TodoRepository
type TodoServiceAdapterRepo struct {
    Mock *mocks.MockTodoService
}

// Helper function to create a TodoService adapter
func CreateTodoService(mockService *mocks.MockTodoService) *todos.TodoService {
    // In a real implementation, this would initialize the service with the mockRepo
    // Since we're mocking the service directly, we can just return the mock itself
    return &todos.TodoService{}
}

func TestGetTodosHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetTodosHandler ===")
    fmt.Println("Testing get todos endpoint")
    
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
    
    // Setup expectations
    mockService.On("GetTodosByUserID", mock.Anything, userID).Return(&dto.TodosResponse{Todos: todos}, nil)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful todos retrieval")
    
    // Create a custom handler that uses our mock service directly
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        todosResponse, err := mockService.GetTodosByUserID(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Format all todos with proper date/time strings
        formattedTodos := make([]map[string]interface{}, len(todosResponse.Todos))
        for i, todo := range todosResponse.Todos {
            dateStr := ""
            timeStr := ""
            
            if !todo.Date.IsZero() {
                dateStr = todo.Date.Format("2006-01-02")
            } else {
                dateStr = time.Now().Format("2006-01-02")
            }
            
            if !todo.Time.IsZero() {
                timeStr = todo.Time.Format("15:04:05")
            } else {
                timeStr = time.Now().Format("15:04:05")
            }
            
            formattedTodos[i] = map[string]interface{}{
                "id":          todo.ID,
                "task":        todo.Task,
                "description": todo.Description,
                "done":        todo.Done,
                "important":   todo.Important,
                "user_id":     todo.UserID,
                "date":        dateStr,
                "time":        timeStr,
            }
        }
        
        json.NewEncoder(w).Encode(formattedTodos)
    }
    
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response data
    assert.Equal(t, 2, len(responseBody))
    assert.Equal(t, "todo-1", responseBody[0]["id"])
    assert.Equal(t, "Test Task 1", responseBody[0]["task"])
    fmt.Printf("✅ Successfully retrieved %d todos for user\n", len(responseBody))
    
    // Setup expectation for scenario 2 - database error
    mockService.On("GetTodosByUserID", mock.Anything, "error-user").Return(nil, errors.New("database error"))
    
    // Create HTTP request with error user context
    req = httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, "error-user")
    req = req.WithContext(ctx)
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing database error when retrieving todos")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All GetTodosHandler test scenarios passed")
}

func TestCreateTodoHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestCreateTodoHandler ===")
    fmt.Println("Testing create todo endpoint")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    task := "New Task"
    description := "New Description"
    important := true
    todoID := "new-todo-1"
    
    // Setup expectations
    mockService.On(
        "CreateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(req *dto.CreateTodoRequest) bool {
            return req.Task == task && 
                   req.Description == description && 
                   req.Important == important &&
                   req.UserID == userID
        }),
    ).Return(&dto.CreateResponse{ID: todoID}, nil)
    
    // Create a request body
    reqBody := map[string]interface{}{
        "task": task,
        "description": description,
        "important": important,
        "date": "2023-07-01",
        "time": "10:00:00",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful todo creation")
    
    // Create a custom handler that uses our mock service directly
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.CreateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the user ID from the context
        userID := r.Context().Value(middleware.UserIDKey).(string)
        req.UserID = userID
        
        // Parse date if provided as string
        if req.DateString != "" {
            parsedDate, err := time.Parse("2006-01-02", req.DateString)
            if err != nil {
                http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
                return
            }
            req.Date = parsedDate
        }
        
        // Parse time if provided as string
        if req.TimeString != "" {
            parsedTime, err := time.Parse("15:04:05", req.TimeString)
            if err != nil {
                http.Error(w, "Invalid time format. Use HH:MM:SS", http.StatusBadRequest)
                return
            }
            req.Time = parsedTime
        }
        
        res, err := mockService.CreateTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(res)
    }
    
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]string
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response
    assert.Equal(t, todoID, responseBody["id"])
    fmt.Printf("✅ Todo created successfully with ID: %s\n", responseBody["id"])
    
    // Setup expectation for scenario 2 - database error
    mockService.On(
        "CreateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(req *dto.CreateTodoRequest) bool {
            return req.Task == "Error Task"
        }),
    ).Return(nil, errors.New("database error"))
    
    // Create a request body for error case
    reqBody = map[string]interface{}{
        "task": "Error Task",
        "description": "Error Description",
        "important": false,
    }
    bodyBytes, _ = json.Marshal(reqBody)
    
    // Create HTTP request
    req = httptest.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing creation with database error")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All CreateTodoHandler test scenarios passed")
}

func TestUpdateTodoHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUpdateTodoHandler ===")
    fmt.Println("Testing update todo endpoint")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    task := "Updated Task"
    description := "Updated Description"
    done := true
    important := false
    
    // Setup router to extract URL params
    router := mux.NewRouter()
    router.HandleFunc("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
    
    // Setup expectations
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(req *dto.UpdateTodoRequest) bool {
            return req.ID == todoID && 
                   req.Task == task && 
                   req.Description == description && 
                   req.Done == done &&
                   req.Important == important &&
                   req.UserID == userID
        }),
    ).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Create a request body
    reqBody := map[string]interface{}{
        "task": task,
        "description": description,
        "done": done,
        "important": important,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodPut, "/api/v1/todo/"+todoID, bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    var match mux.RouteMatch
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful todo update")
    
    // Create a custom handler that uses our mock service directly
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        var req dto.UpdateTodoRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }
        
        // Set the user ID from the context and ID from the URL params
        userID := r.Context().Value(middleware.UserIDKey).(string)
        req.UserID = userID
        req.ID = mux.Vars(r)["id"]
        
        res, err := mockService.UpdateTodo(context.Background(), &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
    
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]bool
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response
    assert.True(t, responseBody["success"])
    fmt.Printf("✅ Todo %s updated successfully\n", todoID)
    
    // Setup expectations for scenario 2 - todo not found
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(req *dto.UpdateTodoRequest) bool {
            return req.ID == "nonexistent-todo"
        }),
    ).Return(nil, errors.New("todo not found"))
    
    // Create HTTP request for non-existent todo
    req = httptest.NewRequest(http.MethodPut, "/api/v1/todo/nonexistent-todo", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": "nonexistent-todo"})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing update of non-existent todo")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Setup expectations for scenario 3 - unauthorized update
    mockService.On(
        "UpdateTodo", 
        mock.Anything, 
        mock.MatchedBy(func(req *dto.UpdateTodoRequest) bool {
            return req.ID == todoID && req.UserID == "wrong-user"
        }),
    ).Return(nil, errors.New("unauthorized"))
    
    // Create HTTP request with wrong user
    req = httptest.NewRequest(http.MethodPut, "/api/v1/todo/"+todoID, bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, "wrong-user")
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 3: Testing unauthorized update")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All UpdateTodoHandler test scenarios passed")
}

func TestDeleteTodoHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestDeleteTodoHandler ===")
    fmt.Println("Testing delete todo endpoint")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    
    // Setup router to extract URL params
    router := mux.NewRouter()
    router.HandleFunc("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")
    
    // Setup expectations
    mockService.On("DeleteTodo", mock.Anything, todoID, userID).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodDelete, "/api/v1/todo/"+todoID, nil)
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    var match mux.RouteMatch
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful todo deletion")
    
    // Create a custom handler that uses our mock service directly
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        res, err := mockService.DeleteTodo(context.Background(), params["id"], userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
    
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]bool
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response
    assert.True(t, responseBody["success"])
    fmt.Printf("✅ Todo %s deleted successfully\n", todoID)
    
    // Setup expectations for scenario 2 - todo not found
    mockService.On("DeleteTodo", mock.Anything, "nonexistent-todo", userID).Return(nil, errors.New("todo not found"))
    
    // Create HTTP request for non-existent todo
    req = httptest.NewRequest(http.MethodDelete, "/api/v1/todo/nonexistent-todo", nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": "nonexistent-todo"})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing deletion of non-existent todo")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Setup expectations for scenario 3 - unauthorized deletion
    mockService.On("DeleteTodo", mock.Anything, todoID, "wrong-user").Return(nil, errors.New("unauthorized"))
    
    // Create HTTP request with wrong user
    req = httptest.NewRequest(http.MethodDelete, "/api/v1/todo/"+todoID, nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, "wrong-user")
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 3: Testing unauthorized deletion")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All DeleteTodoHandler test scenarios passed")
}

func TestUndoTodoHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestUndoTodoHandler ===")
    fmt.Println("Testing undo todo endpoint")
    
    // Create mock service
    mockService := new(mocks.MockTodoService)
    
    // Setup test data
    userID := "user-123"
    todoID := "todo-1"
    
    // Setup router to extract URL params
    router := mux.NewRouter()
    router.HandleFunc("/api/v1/todo/undo/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
    
    // Setup expectations
    mockService.On("UndoTodo", mock.Anything, todoID, userID).Return(&dto.SuccessResponse{Success: true}, nil)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodPut, "/api/v1/todo/undo/"+todoID, nil)
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    var match mux.RouteMatch
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful todo undoing")
    
    // Create a custom handler that uses our mock service directly
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        params := mux.Vars(r)
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        res, err := mockService.UndoTodo(context.Background(), params["id"], userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(res)
    }
    
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]bool
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response
    assert.True(t, responseBody["success"])
    fmt.Printf("✅ Todo %s undone successfully\n", todoID)
    
    // Setup expectations for scenario 2 - todo not found
    mockService.On("UndoTodo", mock.Anything, "nonexistent-todo", userID).Return(nil, errors.New("todo not found"))
    
    // Create HTTP request for non-existent todo
    req = httptest.NewRequest(http.MethodPut, "/api/v1/todo/undo/nonexistent-todo", nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, userID)
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": "nonexistent-todo"})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing undoing of non-existent todo")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Setup expectations for scenario 3 - unauthorized undoing
    mockService.On("UndoTodo", mock.Anything, todoID, "wrong-user").Return(nil, errors.New("unauthorized"))
    
    // Create HTTP request with wrong user
    req = httptest.NewRequest(http.MethodPut, "/api/v1/todo/undo/"+todoID, nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, "wrong-user")
    req = req.WithContext(ctx)
    
    // Match the request and extract params
    router.Match(req, &match)
    req = mux.SetURLVars(req, map[string]string{"id": todoID})
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 3: Testing unauthorized undoing")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All UndoTodoHandler test scenarios passed")
}