package team_todos_repository

import (
    "context"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

type TeamTodoServiceRepository interface {
    CreateTeamTodo(ctx context.Context, task, description string, done, important bool, teamID, assignedTo string) error
    GetTeamTodos(ctx context.Context, teamID string) ([]domain.TeamTodo, error)
    UpdateTeamTodo(ctx context.Context, id, task, description string, done, important bool, teamID, assignedTo string) error
    DeleteTeamTodo(ctx context.Context, id, teamID string) error
}

type TeamTodoRepository struct {
    querier *db.Queries
}

func (r *TeamTodoRepository) CreateTeamTodo(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTeamTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TeamTodoRepository) GetTeamTodos(ctx context.Context, teamID string) (*dto.TeamTodosResponse, error) {
    todos, err := r.querier.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamTodosResponse(todos), nil
}

func (r *TeamTodoRepository) UpdateTeamTodo(ctx context.Context, req *dto.UpdateTeamTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTeamTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TeamTodoRepository) DeleteTeamTodo(ctx context.Context, id, teamID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTeamTodo(ctx, db.DeleteTeamTodoParams{
        ID:     id,
        TeamID: teamID,
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}