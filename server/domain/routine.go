package domain

import (
    "context"
    "time"
)

// Routine represents a recurring task schedule
type Routine struct {
    ID           string    `json:"id"`
    Day          string    `json:"day"`
    ScheduleType string    `json:"scheduleType"`
    TaskID       string    `json:"taskId"`
    UserID       string    `json:"userId"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    IsActive     bool      `json:"isActive"`
}

// RoutineRepository defines the interface for routine persistence operations
type RoutineRepository interface {
    CreateRoutine(ctx context.Context, day, scheduleType, taskID, userID string, isActive bool) (string, error)
    UpdateRoutineStatus(ctx context.Context, id string, isActive bool) error
    UpdateRoutineDay(ctx context.Context, id, day string) error
    GetRoutinesByTaskID(ctx context.Context, taskID string) ([]Routine, error)
    GetDailyRoutines(ctx context.Context, day, scheduleType, userID string) ([]Todo, error)
    DeleteRoutinesByTaskID(ctx context.Context, taskID string) error
    CreateOrUpdateRoutines(ctx context.Context, taskID string, schedules []string, day string, userID string) ([]Routine, error)

}