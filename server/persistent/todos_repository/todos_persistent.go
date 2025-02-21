package todos_repository

import (
    "context"
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type todoRepository struct {
    querier *db.Queries
}


func (r *todoRepository) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *todoRepository) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewTodosResponse(todos), nil
}

func (r *todoRepository) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *todoRepository) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}