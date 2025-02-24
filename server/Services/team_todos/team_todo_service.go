package team_todos

import (
    "context"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_todos_repository"
)

type TeamTodoService struct {
    repo *team_todos_repository.TeamTodoRepository
}

func (s *TeamTodoService) CreateTeamTodo(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.CreateTeamTodo"
    res, err := s.repo.CreateTeamTodo(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create team todo: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamTodoService) GetTeamTodos(ctx context.Context, teamID string) (*dto.TeamTodosResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.GetTeamTodos"
    res, err := s.repo.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team todos: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamTodoService) UpdateTeamTodo(ctx context.Context, req *dto.UpdateTeamTodoRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.UpdateTeamTodo"
    res, err := s.repo.UpdateTeamTodo(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update team todo: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamTodoService) DeleteTeamTodo(ctx context.Context, id, teamID string) (*dto.SuccessResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.DeleteTeamTodo"
    res, err := s.repo.DeleteTeamTodo(ctx, id, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete team todo: %w", functionName, err)
    }
    return res, nil
}