package todos_repository

import (
    "context"
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "time"
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
    
    // Convert db.Todo to domain.Todo with proper time handling
    domainTodos := make([]domain.Todo, len(todos))
    for i, todo := range todos {
        var date time.Time
        var timeVal time.Time
        
        // Handle date conversion safely
        if todo.Date.Valid {
            date = todo.Date.Time
        }
        
        // Handle time conversion safely
        if todo.Time.Valid {
            timeVal = todo.Time.Time
        }
        
        domainTodos[i] = domain.Todo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            Date:        date,
            Time:        timeVal,
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
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewTodosResponse(todos), nil
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