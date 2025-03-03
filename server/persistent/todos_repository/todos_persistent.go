package todos_repository

import (
    "context"
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "time"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure TodoRepository implements domain.TodoRepository
var _ domain.TodoRepository = (*TodoRepository)(nil)

type TodoRepository struct {
    querier *db.Queries
}

// Implement domain.TodoRepository interface methods
func (r *TodoRepository) CreateTodo(ctx context.Context, task, description string, done, important bool, userID string) (string, error) {
    // Use your existing DTO and converter
    req := &dto.CreateTodoRequest{
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    params := req.ConvertCreateTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTodo(ctx, *params)
    if err != nil {
        return "", err
    }
    return params.ID, nil
}

func (r *TodoRepository) GetTodosByUserID(ctx context.Context, userID string) ([]domain.Todo, error) {
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Log the types to understand what we're dealing with
    if len(todos) > 0 {
        fmt.Printf("Date type: %T\n", todos[0].Date)
        fmt.Printf("Time type: %T\n", todos[0].Time)
    }
    
    // Convert db.Todo to domain.Todo
    domainTodos := make([]domain.Todo, len(todos))
    for i, todo := range todos {
        domainTodos[i] = domain.Todo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            // Skip date/time for now to see if code works otherwise
        }
        
        // Add a type switch to handle different date/time types
        switch d := todo.Date.(type) {
        case time.Time:
            domainTodos[i].Date = d
        case string:
            parsed, _ := time.Parse("2006-01-02", d)
            domainTodos[i].Date = parsed
        }
        
        switch t := todo.Time.(type) {
        case time.Time:
            domainTodos[i].Time = t
        case string:
            parsed, _ := time.Parse("15:04:05", t)
            domainTodos[i].Time = parsed
        }
    }
    
    return domainTodos, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) (bool, error) {
    // Use your existing DTO and converter
    req := &dto.UpdateTodoRequest{
        ID:          id,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id, userID string) (bool, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TodoRepository) UndoTodo(ctx context.Context, id, userID string) (bool, error) {
    err := r.querier.UndoTodo(ctx, db.UndoTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

// Original methods for backward compatibility - you can keep these if you want
func (r *TodoRepository) CreateTodoWithDTO(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TodoRepository) GetTodosByUserIDWithDTO(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    todosRows, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Convert []db.GetTodosByUserIDRow to []dto.TodoResponse
    var todoResponses []dto.TodoResponse
    for _, todo := range todosRows {
        // Handle date and time based on their types
        var date time.Time
        var timeVal time.Time
        
        // Handle date conversion
        switch d := todo.Date.(type) {
        case time.Time:
            date = d
        case string:
            parsed, _ := time.Parse("2006-01-02", d)
            date = parsed
        }
        
        // Handle time conversion
        switch t := todo.Time.(type) {
        case time.Time:
            timeVal = t
        case string:
            parsed, _ := time.Parse("15:04:05", t)
            timeVal = parsed
        }
        
        todoResponses = append(todoResponses, dto.TodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            Date:        date,
            Time:        timeVal,
        })
    }
    
    return &dto.TodosResponse{Todos: todoResponses}, nil
}

func (r *TodoRepository) UpdateTodoWithDTO(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) DeleteTodoWithDTO(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) UndoTodoWithDTO(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.UndoTodo(ctx, db.UndoTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}