package shared_todos

import (
    "context"
    "fmt"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type SharedTodoService struct {
    repo domain.SharedTodoRepository
    todoRepo domain.TodoRepository
    userRepo domain.UserRepository
}


// ShareTodo shares a todo with another user
func (s *SharedTodoService) ShareTodo(ctx context.Context, todoID string, recipientUserID string, sharedBy string) error {
    const functionName = "services.shared_todos.SharedTodoService.ShareTodo"
    
    // Get the original todo to make sure it exists and belongs to the current user
    todo, err := s.todoRepo.GetTodoByID(ctx, todoID)
    if err != nil {
        return fmt.Errorf("%s: failed to get todo: %w", functionName, err)
    }
    
    // Verify the todo belongs to the current user
    if todo.UserID != sharedBy {
        return fmt.Errorf("%s: unauthorized to share this todo", functionName)
    }
    
    // Check if the todo is already shared with this user
    isShared, err := s.repo.IsSharedWithUser(ctx, todoID, recipientUserID)
    if err != nil {
        return fmt.Errorf("%s: failed to check if todo is shared: %w", functionName, err)
    }
    
    if isShared {
        return fmt.Errorf("%s: todo is already shared with this user", functionName)
    }
    
    // Share the todo
    err = s.repo.ShareTodo(ctx, todoID, recipientUserID, sharedBy)
    if err != nil {
        return fmt.Errorf("%s: failed to share todo: %w", functionName, err)
    }
    
    return nil
}

// In server/services/shared_todos/shared_todo_service.go
func (s *SharedTodoService) CreateSharedTodo(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.CreateSharedTodo"
    
    // Validate required fields
    if req.Task == "" {
        return nil, fmt.Errorf("%s: task cannot be empty", functionName)
    }
    
    // Handle date value - if your CreateSharedTodoRequest has date fields
    date := req.Date
    if date.IsZero() {
        // If DateString is provided, parse it
        if req.DateString != "" {
            parsedDate, err := time.Parse("2006-01-02", req.DateString)
            if err == nil {
                date = parsedDate
            } else {
                // If parsing fails, default to today
                date = time.Now()
            }
        } else {
            // If no date provided at all, default to today
            date = time.Now()
        }
    }
    
    // Handle time value - if your CreateSharedTodoRequest has time fields
    timeValue := req.Time
    if timeValue.IsZero() {
        if req.TimeString != "" {
            parsedTime, err := time.Parse("15:04:05", req.TimeString)
            if err == nil {
                // Extract just the time parts and use a valid year
                hour, min, sec := parsedTime.Clock()
                timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            } else {
                // If parsing fails, default to current time
                now := time.Now()
                timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            }
        } else {
            // If no time provided, default to current time
            now := time.Now()
            timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        }
    } else if timeValue.Year() < 1 {
        // Ensure the year is valid
        hour, min, sec := timeValue.Clock()
        timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
    }
    
    id, err := s.repo.CreateSharedTodo(ctx, req.Task, req.Description, req.Done, req.Important, req.UserID, req.SharedBy)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create shared todo: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *SharedTodoService) GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedTodos"
    domainTodos, err := s.repo.GetSharedTodos(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared todos: %w", functionName, err)
    }
    
    // Convert domain.SharedTodo to dto.SharedTodoResponse
    var receivedTodos []dto.SharedTodoResponse
    for _, todo := range domainTodos {
        receivedTodos = append(receivedTodos, dto.SharedTodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
            SharedBy:    todo.SharedBy,
        })
    }
    
    return &dto.SharedTodosResponse{Received: receivedTodos}, nil
}

func (s *SharedTodoService) GetSharedByMeTodos(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedByMeTodos"
    domainTodos, err := s.repo.GetSharedByMeTodos(ctx, sharedBy)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared by me todos: %w", functionName, err)
    }
    
    // Convert domain.SharedTodo to dto.SharedTodoResponse
    var sharedTodos []dto.SharedTodoResponse
    for _, todo := range domainTodos {
        sharedTodos = append(sharedTodos, dto.SharedTodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
            SharedBy:    todo.SharedBy,
        })
    }
    
    return &dto.SharedTodosResponse{Shared: sharedTodos}, nil
}

