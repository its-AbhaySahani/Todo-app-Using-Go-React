package team_todos_repository

import (
    "context"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure TeamTodoRepository implements domain.TeamTodoRepository
var _ domain.TeamTodoRepository = (*TeamTodoRepository)(nil)

type TeamTodoRepository struct {
    querier *db.Queries
}

// Implement domain.TeamTodoRepository interface methods
func (r *TeamTodoRepository) CreateTeamTodo(ctx context.Context, task, description string, done, important bool, teamID, assignedTo string) (string, error) {
    // Use your existing DTO and converter
    req := &dto.CreateTeamTodoRequest{
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        TeamID:      teamID,
        AssignedTo:  assignedTo,
    }
    
    params := req.ConvertCreateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTeamTodo(ctx, *params)
    if err != nil {
        return "", err
    }
    return params.ID, nil
}

func (r *TeamTodoRepository) GetTeamTodos(ctx context.Context, teamID string) ([]domain.TeamTodo, error) {
    todos, err := r.querier.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, err
    }
    
    // Convert db.TeamTodo to domain.TeamTodo
    domainTodos := make([]domain.TeamTodo, len(todos))
    for i, todo := range todos {
        domainTodos[i] = domain.TeamTodo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important.Bool,
            TeamID:      todo.TeamID,
            AssignedTo:  todo.AssignedTo.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
        }
    }
    
    return domainTodos, nil
}

func (r *TeamTodoRepository) UpdateTeamTodo(ctx context.Context, id, task, description string, done, important bool, teamID, assignedTo string) (bool, error) {
    // Use your existing DTO and converter
    req := &dto.UpdateTeamTodoRequest{
        ID:          id,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        TeamID:      teamID,
        AssignedTo:  assignedTo,
    }
    
    params := req.ConvertUpdateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTeamTodo(ctx, *params)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TeamTodoRepository) DeleteTeamTodo(ctx context.Context, id, teamID string) (bool, error) {
    err := r.querier.DeleteTeamTodo(ctx, db.DeleteTeamTodoParams{
        ID:     id,
        TeamID: teamID,
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

// Original methods for backward compatibility
func (r *TeamTodoRepository) CreateTeamTodoWithDTO(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTeamTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TeamTodoRepository) GetTeamTodosWithDTO(ctx context.Context, teamID string) (*dto.TeamTodosResponse, error) {
    todos, err := r.querier.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamTodosResponse(todos), nil
}

func (r *TeamTodoRepository) UpdateTeamTodoWithDTO(ctx context.Context, req *dto.UpdateTeamTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTeamTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTeamTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TeamTodoRepository) DeleteTeamTodoWithDTO(ctx context.Context, id, teamID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTeamTodo(ctx, db.DeleteTeamTodoParams{
        ID:     id,
        TeamID: teamID,
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}