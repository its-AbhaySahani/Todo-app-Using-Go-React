package team_todos

import (
    "context"
    "fmt"
    "time"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TeamTodoService struct {
    repo domain.TeamTodoRepository
}

// In server/services/team_todos/team_todo_service.go
func (s *TeamTodoService) CreateTeamTodo(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.team_todos.TeamTodoService.CreateTeamTodo"
    
    // Validate required fields
    if req.Task == "" {
        return nil, fmt.Errorf("%s: task cannot be empty", functionName)
    }
    
    // Handle date value - if your CreateTeamTodoRequest has date fields
    date := req.Date
    if date.IsZero() {
        // If DateString is provided, parse it
        if req.DateString != "" {
            parsedDate, err := time.Parse("2006-01-02", req.DateString)
            if err == nil {
                date = parsedDate
            } else {
                // If parsing fails, default to today
                date = time.Now()
            }
        } else {
            // If no date provided at all, default to today
            date = time.Now()
        }
    }
    
    // Handle time value - if your CreateTeamTodoRequest has time fields
    timeValue := req.Time
    if timeValue.IsZero() {
        if req.TimeString != "" {
            parsedTime, err := time.Parse("15:04:05", req.TimeString)
            if err == nil {
                // Extract just the time parts and use a valid year
                hour, min, sec := parsedTime.Clock()
                timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            } else {
                // If parsing fails, default to current time
                now := time.Now()
                timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            }
        } else {
            // If no time provided, default to current time
            now := time.Now()
            timeValue = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        }
    } else if timeValue.Year() < 1 {
        // Ensure the year is valid
        hour, min, sec := timeValue.Clock()
        timeValue = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
    }
    
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