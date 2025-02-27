package team_todos

import (
    "context"
    "fmt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TeamTodoService struct {
    repo domain.TeamTodoRepository
}

func (s *TeamTodoService) CreateTeamTodo(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.CreateTeamTodo"
    id, err := s.repo.CreateTeamTodo(ctx, req.Task, req.Description, req.Done, req.Important, req.TeamID, req.AssignedTo)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create team todo: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *TeamTodoService) GetTeamTodos(ctx context.Context, teamID string) (*dto.TeamTodosResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.GetTeamTodos"
    domainTodos, err := s.repo.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team todos: %w", functionName, err)
    }
    
    // Convert domain.TeamTodo to dto.TeamTodoResponse
    var todoResponses []dto.TeamTodoResponse
    for _, todo := range domainTodos {
        todoResponses = append(todoResponses, dto.TeamTodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            TeamID:      todo.TeamID,
            AssignedTo:  todo.AssignedTo,
            Date:        todo.Date,
            Time:        todo.Time,
        })
    }
    
    return &dto.TeamTodosResponse{Todos: todoResponses}, nil
}

func (s *TeamTodoService) UpdateTeamTodo(ctx context.Context, req *dto.UpdateTeamTodoRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.UpdateTeamTodo"
    success, err := s.repo.UpdateTeamTodo(ctx, req.ID, req.Task, req.Description, req.Done, req.Important, req.TeamID, req.AssignedTo)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update team todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TeamTodoService) DeleteTeamTodo(ctx context.Context, id, teamID string) (*dto.SuccessResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.DeleteTeamTodo"
    success, err := s.repo.DeleteTeamTodo(ctx, id, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete team todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}