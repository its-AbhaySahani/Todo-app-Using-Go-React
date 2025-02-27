package domain

import (
    "context"
)

type User struct {
    ID       string
    Username string
    Password string
}

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
    CreateUser(ctx context.Context, username, password string) (string, error)
    GetUserByUsername(ctx context.Context, username string) (User, error)
}