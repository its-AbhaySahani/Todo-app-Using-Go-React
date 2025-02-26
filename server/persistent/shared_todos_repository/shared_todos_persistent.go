package shared_todos_repository

import (
	"context"
	"database/sql"

	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type SharedTodoServiceRepository interface {
	CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) error
	GetSharedTodos(ctx context.Context, userID string) ([]domain.SharedTodo, error)
	GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]domain.SharedTodo, error)
}

type SharedTodoRepository struct {
	querier *db.Queries
}

func (r *SharedTodoRepository) CreateSharedTodo(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateSharedTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateSharedTodo(ctx, *params)      // Using sqlc.arg
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *SharedTodoRepository) GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedTodos(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}

func (r *SharedTodoRepository) GetSharedByMeTodos(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedByMeTodos(ctx, sql.NullString{String: sharedBy, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}
