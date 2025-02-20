package domain

import "context"

type User struct {
    ID       string
    Username string
    Password string
}

type UserRepository interface {
    CreateUser(ctx context.Context, username, password string) error
    GetUserByUsername(ctx context.Context, username string) (User, error)
}