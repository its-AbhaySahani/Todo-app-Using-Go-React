package shared_todos_repository

import (
    "context"
    "database/sql"
    "time"
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure SharedTodoRepository implements domain.SharedTodoRepository
var _ domain.SharedTodoRepository = (*SharedTodoRepository)(nil)

type SharedTodoRepository struct {
    querier *db.Queries
}

// Implement domain.SharedTodoRepository interface methods
func (r *SharedTodoRepository) CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) (string, error) {
    id := uuid.New().String()
    
    // Use current date and time if none provided
    nullDate := sql.NullTime{Time: time.Now(), Valid: true}
    nullTime := sql.NullTime{Time: time.Now(), Valid: true}
    
    err := r.querier.CreateSharedTodo(ctx, db.CreateSharedTodoParams{
        ID:          id,
        Task:        sql.NullString{String: task, Valid: true},
        Description: sql.NullString{String: description, Valid: true},
        Done:        sql.NullBool{Bool: done, Valid: true},
        Important:   sql.NullBool{Bool: important, Valid: true},
        UserID:      sql.NullString{String: userID, Valid: true},
        Date:        nullDate,
        Time:        nullTime,
        SharedBy:    sql.NullString{String: sharedBy, Valid: true},
    })
    
    if err != nil {
        return "", err
    }
    
    return id, nil
}

func (r *SharedTodoRepository) GetSharedTodos(ctx context.Context, userID string) ([]domain.SharedTodo, error) {
    todos, err := r.querier.GetSharedTodos(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Convert db.SharedTodo to domain.SharedTodo
    domainTodos := make([]domain.SharedTodo, len(todos))
    for i, todo := range todos {
        domainTodos[i] = domain.SharedTodo{
            ID:          todo.ID,
            Task:        todo.Task.String,
            Description: todo.Description.String,
            Done:        todo.Done.Bool,
            Important:   todo.Important.Bool,
            UserID:      todo.UserID.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
            SharedBy:    todo.SharedBy.String,
        }
    }
    
    return domainTodos, nil
}

func (r *SharedTodoRepository) GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]domain.SharedTodo, error) {
    todos, err := r.querier.GetSharedByMeTodos(ctx, sql.NullString{String: sharedBy, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Convert db.SharedTodo to domain.SharedTodo
    domainTodos := make([]domain.SharedTodo, len(todos))
    for i, todo := range todos {
        domainTodos[i] = domain.SharedTodo{
            ID:          todo.ID,
            Task:        todo.Task.String,
            Description: todo.Description.String,
            Done:        todo.Done.Bool,
            Important:   todo.Important.Bool,
            UserID:      todo.UserID.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
            SharedBy:    todo.SharedBy.String,
        }
    }
    
    return domainTodos, nil
}

// Original methods for backward compatibility
func (r *SharedTodoRepository) CreateSharedTodoWithDTO(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateSharedTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateSharedTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *SharedTodoRepository) GetSharedTodosWithDTO(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedTodos(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}

func (r *SharedTodoRepository) GetSharedByMeTodosWithDTO(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedByMeTodos(ctx, sql.NullString{String: sharedBy, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}

