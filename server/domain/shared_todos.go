package domain

import (
    "context"
    "time"
)

// SharedTodo entity
type SharedTodo struct {
    ID          string
    Task        string
    Description string
    Done        bool
    Important   bool
    UserID      string
    Date        time.Time
    Time        time.Time
    SharedBy    string
}

// SharedTodoRepository defines the interface for shared todo persistence operations
type SharedTodoRepository interface {
    CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) (string, error)
    GetSharedTodos(ctx context.Context, userID string) ([]SharedTodo, error)
    GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]SharedTodo, error)
}