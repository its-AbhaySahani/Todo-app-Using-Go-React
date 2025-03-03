package routines

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type RoutineService struct {
    repo domain.RoutineRepository
}

// CreateRoutine creates a new routine
func (s *RoutineService) CreateRoutine(ctx context.Context, req *dto.CreateRoutineRequest) (*dto.CreateResponse, error) {
    const functionName = "services.routines.RoutineService.CreateRoutine"
    
    id, err := s.repo.CreateRoutine(ctx, req.Day, req.ScheduleType, req.TaskID, req.UserID, req.IsActive)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create routine: %w", functionName, err)
    }
    
    return &dto.CreateResponse{ID: id}, nil
}

// UpdateRoutineStatus updates the active status of a routine
func (s *RoutineService) UpdateRoutineStatus(ctx context.Context, id string, isActive bool) (*dto.SuccessResponse, error) {
    const functionName = "services.routines.RoutineService.UpdateRoutineStatus"
    
    err := s.repo.UpdateRoutineStatus(ctx, id, isActive)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update routine status: %w", functionName, err)
    }
    
    return &dto.SuccessResponse{Success: true}, nil
}

// UpdateRoutineDay updates the day of a routine
func (s *RoutineService) UpdateRoutineDay(ctx context.Context, id, day string) (*dto.SuccessResponse, error) {
    const functionName = "services.routines.RoutineService.UpdateRoutineDay"
    
    err := s.repo.UpdateRoutineDay(ctx, id, day)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update routine day: %w", functionName, err)
    }
    
    return &dto.SuccessResponse{Success: true}, nil
}

// GetRoutinesByTaskID gets routines for a specific task
func (s *RoutineService) GetRoutinesByTaskID(ctx context.Context, taskID string) (*dto.RoutinesResponse, error) {
    const functionName = "services.routines.RoutineService.GetRoutinesByTaskID"
    
    routines, err := s.repo.GetRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get routines: %w", functionName, err)
    }
    
    // Convert domain routines to DTO response
    var routineResponses []dto.RoutineResponse
    for _, routine := range routines {
        routineResponses = append(routineResponses, dto.RoutineResponse{
            ID:           routine.ID,
            Day:          routine.Day,
            ScheduleType: routine.ScheduleType,
            TaskID:       routine.TaskID,
            UserID:       routine.UserID,
            CreatedAt:    routine.CreatedAt,
            UpdatedAt:    routine.UpdatedAt,
            IsActive:     routine.IsActive,
        })
    }
    
    return &dto.RoutinesResponse{Routines: routineResponses}, nil
}

// GetDailyRoutines gets todos for a specific day and schedule type
func (s *RoutineService) GetDailyRoutines(ctx context.Context, day, scheduleType, userID string) (*dto.TodosResponse, error) {
    const functionName = "services.routines.RoutineService.GetDailyRoutines"
    
    todos, err := s.repo.GetDailyRoutines(ctx, day, scheduleType, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get daily routines: %w", functionName, err)
    }
    
    // Convert domain todos to DTO response
    var todoResponses []dto.TodoResponse
    for _, todo := range todos {
        todoResponses = append(todoResponses, dto.TodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
        })
    }
    
    return &dto.TodosResponse{Todos: todoResponses}, nil
}

// GetTodayRoutines gets todos for today's routines by schedule type
func (s *RoutineService) GetTodayRoutines(ctx context.Context, scheduleType, userID string) (*dto.TodosResponse, error) {
    const functionName = "services.routines.RoutineService.GetTodayRoutines"
    
    // Get today's day name (sunday, monday, etc.)
    dayName := strings.ToLower(time.Now().Weekday().String())
    
    return s.GetDailyRoutines(ctx, dayName, scheduleType, userID)
}

// DeleteRoutinesByTaskID deletes all routines for a task
func (s *RoutineService) DeleteRoutinesByTaskID(ctx context.Context, taskID string) (*dto.SuccessResponse, error) {
    const functionName = "services.routines.RoutineService.DeleteRoutinesByTaskID"
    
    err := s.repo.DeleteRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete routines: %w", functionName, err)
    }
    
    return &dto.SuccessResponse{Success: true}, nil
}

// CreateOrUpdateRoutines creates or updates routines for a task
func (s *RoutineService) CreateOrUpdateRoutines(ctx context.Context, req *dto.CreateOrUpdateRoutinesRequest) (*dto.RoutinesResponse, error) {
    const functionName = "services.routines.RoutineService.CreateOrUpdateRoutines"
    
    // Use current day if not provided
    day := req.Day
    if day == "" {
        day = strings.ToLower(time.Now().Weekday().String())
    }
    
    // Call the repository method
    routines, err := s.repo.CreateOrUpdateRoutines(ctx, req.TaskID, req.Schedules, day, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create or update routines: %w", functionName, err)
    }
    
    // Convert domain routines to DTO routines
    var routineResponses []dto.RoutineResponse
    for _, routine := range routines {
        routineResponses = append(routineResponses, dto.RoutineResponse{
            ID:           routine.ID,
            Day:          routine.Day,
            ScheduleType: routine.ScheduleType,
            TaskID:       routine.TaskID,
            UserID:       routine.UserID,
            CreatedAt:    routine.CreatedAt,
            UpdatedAt:    routine.UpdatedAt,
            IsActive:     routine.IsActive,
        })
    }
    
    return &dto.RoutinesResponse{Routines: routineResponses}, nil
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}