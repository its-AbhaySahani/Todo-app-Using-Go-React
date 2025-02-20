package domain

import (
    "context"
    "time"
)

type Todo struct {
    ID          string
    Task        string
    Description string
    Done        bool
    Important   bool
    UserID      string
    Date        time.Time
    Time        time.Time
}

type TodoRepository interface {
    CreateTodo(ctx context.Context, task, description string, done, important bool, userID string) error
    GetTodosByUserID(ctx context.Context, userID string) ([]Todo, error)
    UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) error
    DeleteTodo(ctx context.Context, id, userID string) error
}