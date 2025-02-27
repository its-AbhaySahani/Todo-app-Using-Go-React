package domain

import (
    "context"
    "time"
)

type TeamTodo struct {
    ID          string
    Task        string
    Description string
    Done        bool
    Important   bool
    TeamID      string
    AssignedTo  string
    Date        time.Time
    Time        time.Time
}

// TeamTodoRepository defines the interface for team todo persistence operations
type TeamTodoRepository interface {
    CreateTeamTodo(ctx context.Context, task, description string, done, important bool, teamID, assignedTo string) (string, error)
    GetTeamTodos(ctx context.Context, teamID string) ([]TeamTodo, error)
    UpdateTeamTodo(ctx context.Context, id, task, description string, done, important bool, teamID, assignedTo string) (bool, error)
    DeleteTeamTodo(ctx context.Context, id, teamID string) (bool, error)
}