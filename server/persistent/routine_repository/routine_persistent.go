package routine_repository

import (
    "context"
    "time"
    "strings"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure RoutineRepository implements domain.RoutineRepository
var _ domain.RoutineRepository = (*RoutineRepository)(nil)

type RoutineRepository struct {
    querier *db.Queries
}

// Implement domain.RoutineRepository interface methods
func (r *RoutineRepository) CreateRoutine(ctx context.Context, day, scheduleType, taskID, userID string, isActive bool) (string, error) {
    // Use your existing DTO and converter
    req := &dto.CreateRoutineRequest{
        Day:          day,
        ScheduleType: scheduleType,
        TaskID:       taskID,
        UserID:       userID,
        IsActive:     isActive,
    }
    
    params := req.ConvertCreateRoutineDomainRequestToPersistentRequest()
    err := r.querier.CreateRoutine(ctx, *params)
    if err != nil {
        return "", err
    }
    return params.ID, nil
}

func (r *RoutineRepository) UpdateRoutineStatus(ctx context.Context, id string, isActive bool) error {
    req := &dto.UpdateRoutineStatusRequest{
        ID:       id,
        IsActive: isActive,
    }
    
    params := req.ConvertUpdateRoutineStatusDomainRequestToPersistentRequest()
    return r.querier.UpdateRoutineStatus(ctx, *params)
}

func (r *RoutineRepository) UpdateRoutineDay(ctx context.Context, id, day string) error {
    req := &dto.UpdateRoutineDayRequest{
        ID:  id,
        Day: day,
    }
    
    params := req.ConvertUpdateRoutineDayDomainRequestToPersistentRequest()
    return r.querier.UpdateRoutineDay(ctx, *params)
}

func (r *RoutineRepository) GetRoutinesByTaskID(ctx context.Context, taskID string) ([]domain.Routine, error) {
    routines, err := r.querier.GetRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, err
    }
    
    var result []domain.Routine
    for _, routine := range routines {
        result = append(result, domain.Routine{
            ID:           routine.ID,
            Day:          string(routine.Day),
            ScheduleType: string(routine.Scheduletype),
            TaskID:       routine.Taskid,
            UserID:       routine.Userid,
            CreatedAt:    routine.Createdat,
            UpdatedAt:    routine.Updatedat,
            IsActive:     routine.Isactive.Bool,
        })
    }
    
    return result, nil
}

func (r *RoutineRepository) GetDailyRoutines(ctx context.Context, day, scheduleType, userID string) ([]domain.Todo, error) {
    todos, err := r.querier.GetDailyRoutines(ctx, db.GetDailyRoutinesParams{
        Day:          db.RoutinesDay(day),
        Scheduletype: db.RoutinesScheduletype(scheduleType),
        Userid:       userID,
    })
    if err != nil {
        return nil, err
    }
    
    var result []domain.Todo
    for _, todo := range todos {
        result = append(result, domain.Todo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
        })
    }
    
    return result, nil
}

// GetTodayRoutines gets routines for today by schedule type
func (r *RoutineRepository) GetTodayRoutines(ctx context.Context, scheduleType, userID string) ([]domain.Todo, error) {
    day := time.Now().Weekday().String()
    return r.GetDailyRoutines(ctx, day, scheduleType, userID)
}

func (r *RoutineRepository) DeleteRoutinesByTaskID(ctx context.Context, taskID string) error {
    return r.querier.DeleteRoutinesByTaskID(ctx, taskID)
}

// CreateOrUpdateRoutines handles creating or updating routines for a task
func (r *RoutineRepository) CreateOrUpdateRoutines(ctx context.Context, taskID string, schedules []string, day string, userID string) ([]domain.Routine, error) {
    // Use current day if not provided
    if day == "" {
        day = time.Now().Weekday().String()
        day = day[0:1] + strings.ToLower(day[1:]) // Format to lowercase with first letter capital
    }
    
    // Get existing routines for this task
    existingRoutines, err := r.GetRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, err
    }
    
    // Filter routines for this user
    var userRoutines []domain.Routine
    for _, routine := range existingRoutines {
        if routine.UserID == userID {
            userRoutines = append(userRoutines, routine)
        }
    }
    
    // Track which schedule types we process
    scheduleTypesProcessed := make(map[string]bool)
    var updatedRoutines []domain.Routine
    
    // Process existing routines - update or deactivate them
    for _, routine := range userRoutines {
        // Check if this schedule type is in the request
        isInRequest := false
        for _, scheduleType := range schedules {
            if scheduleType == routine.ScheduleType {
                isInRequest = true
                break
            }
        }
        
        if isInRequest {
            // Schedule type is requested, check if day needs updating
            if routine.Day != day {
                // Update the day
                err := r.UpdateRoutineDay(ctx, routine.ID, day)
                if err != nil {
                    return nil, err
                }
                routine.Day = day
            }
            
            // Make sure the routine is active
            if !routine.IsActive {
                err := r.UpdateRoutineStatus(ctx, routine.ID, true)
                if err != nil {
                    return nil, err
                }
                routine.IsActive = true
            }
            
            // Mark this schedule type as processed
            scheduleTypesProcessed[routine.ScheduleType] = true
            updatedRoutines = append(updatedRoutines, routine)
        } else if routine.Day == day {
            // Not requested for this day, deactivate it
            err := r.UpdateRoutineStatus(ctx, routine.ID, false)
            if err != nil {
                return nil, err
            }
        }
    }
    
    // Create new routines for schedule types not yet processed
    for _, scheduleType := range schedules {
        if !scheduleTypesProcessed[scheduleType] {
            // Create a new routine
            id, err := r.CreateRoutine(ctx, day, scheduleType, taskID, userID, true)
            if err != nil {
                return nil, err
            }
            
            // Add to result
            currentTime := time.Now()
            updatedRoutines = append(updatedRoutines, domain.Routine{
                ID:           id,
                Day:          day,
                ScheduleType: scheduleType,
                TaskID:       taskID,
                UserID:       userID,
                CreatedAt:    currentTime,
                UpdatedAt:    currentTime,
                IsActive:     true,
            })
        }
    }
    
    return updatedRoutines, nil
}

// Original methods for backward compatibility
func (r *RoutineRepository) CreateRoutineWithDTO(ctx context.Context, req *dto.CreateRoutineRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateRoutineDomainRequestToPersistentRequest()
    err := r.querier.CreateRoutine(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *RoutineRepository) UpdateRoutineStatusWithDTO(ctx context.Context, req *dto.UpdateRoutineStatusRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateRoutineStatusDomainRequestToPersistentRequest()
    err := r.querier.UpdateRoutineStatus(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *RoutineRepository) UpdateRoutineDayWithDTO(ctx context.Context, req *dto.UpdateRoutineDayRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateRoutineDayDomainRequestToPersistentRequest()
    err := r.querier.UpdateRoutineDay(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *RoutineRepository) GetRoutinesByTaskIDWithDTO(ctx context.Context, taskID string) (*dto.RoutinesResponse, error) {
    routines, err := r.querier.GetRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, err
    }
    return dto.NewRoutinesResponse(routines), nil
}

func (r *RoutineRepository) GetDailyRoutinesWithDTO(ctx context.Context, day, scheduleType, userID string) (*dto.TodosResponse, error) {
    todos, err := r.querier.GetDailyRoutines(ctx, db.GetDailyRoutinesParams{
        Day:          db.RoutinesDay(day),
        Scheduletype: db.RoutinesScheduletype(scheduleType),
        Userid:       userID,
    })
    if err != nil {
        return nil, err
    }
    
    var todosResponse dto.TodosResponse
    for _, todo := range todos {
        todosResponse.Todos = append(todosResponse.Todos, *dto.NewTodoResponse(&todo))
    }
    
    return &todosResponse, nil
}

func (r *RoutineRepository) GetTodayRoutinesWithDTO(ctx context.Context, scheduleType, userID string) (*dto.TodosResponse, error) {
    day := strings.ToLower(time.Now().Weekday().String())
    return r.GetDailyRoutinesWithDTO(ctx, day, scheduleType, userID)
}

func (r *RoutineRepository) DeleteRoutinesByTaskIDWithDTO(ctx context.Context, taskID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteRoutinesByTaskID(ctx, taskID)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *RoutineRepository) CreateOrUpdateRoutinesWithDTO(ctx context.Context, req *dto.CreateOrUpdateRoutinesRequest) (*dto.RoutinesResponse, error) {
    // Use current day if not provided
    day := req.Day
    if day == "" {
        day = strings.ToLower(time.Now().Weekday().String())
    }
    
    routines, err := r.CreateOrUpdateRoutines(ctx, req.TaskID, req.Schedules, day, req.UserID)
    if err != nil {
        return nil, err
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