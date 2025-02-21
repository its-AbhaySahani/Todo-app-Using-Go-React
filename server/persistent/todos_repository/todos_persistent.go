package todos_repository

import (
    "context"
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

type TodoServiceRepository interface {
    CreateTodo(ctx context.Context, task, description string, done, important bool, userID string) error
    GetTodosByUserID(ctx context.Context, userID string) ([]domain.Todo, error)
    UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) error
    DeleteTodo(ctx context.Context, id, userID string) error
    UndoTodo(ctx context.Context, id, userID string) error
}
type TodoRepository struct {
    querier *db.Queries
}


func (r *TodoRepository) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TodoRepository) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewTodosResponse(todos), nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.UndoTodo(ctx, db.UndoTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}