package todos

import (
    "context"
    "fmt"
    "time"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TodoService struct {
    repo domain.TodoRepository
}


func (s *TodoService) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.todos.TodoService.CreateTodo"
    
    // Validate required fields
    if req.Task == "" {
        return nil, fmt.Errorf("%s: task cannot be empty", functionName)
    }
    
    // Handle date value
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
    
    // Handle time value
    timeValue := req.Time
    if timeValue.IsZero() {
        if req.TimeString != "" {
            parsedTime, err := time.Parse("15:04:05", req.TimeString)
            if err == nil {
                // Extract just the time parts and use a valid year (e.g., 2000)
                hour, min, sec := parsedTime.Clock()
                timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            } else {
                // If parsing fails, default to current time with a valid year
                now := time.Now()
                timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            }
        } else {
            // If no time provided, default to current time with a valid year
            now := time.Now()
            timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        }
    } else if timeValue.Year() < 1 {
        // Ensure the year is valid
        hour, min, sec := timeValue.Clock()
        timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
    }
    
    id, err := s.repo.CreateTodo(ctx, req.Task, req.Description, req.Done, req.Important, req.UserID, date, timeValue)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create todo: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *TodoService) GetTodoByID(ctx context.Context, id string) (*domain.Todo, error) {
    const functionName = "services.todos.TodoService.GetTodoByID"
    
    todo, err := s.repo.GetTodoByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get todo by ID: %w", functionName, err)
    }
    
    return todo, nil
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    const functionName = "services.todos.TodoService.GetTodosByUserID"
    domainTodos, err := s.repo.GetTodosByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get todos by user ID: %w", functionName, err)
    }
    
    // Convert domain.Todo to dto.TodoResponse
    var todoResponses []dto.TodoResponse
    for _, todo := range domainTodos {
        todoResponses = append(todoResponses, dto.TodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
        })
    }
    
    return &dto.TodosResponse{Todos: todoResponses}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UpdateTodo"
    success, err := s.repo.UpdateTodo(ctx, req.ID, req.Task, req.Description, req.Done, req.Important, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.DeleteTodo"
    success, err := s.repo.DeleteTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TodoService) UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UndoTodo"
    success, err := s.repo.UndoTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to undo todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}