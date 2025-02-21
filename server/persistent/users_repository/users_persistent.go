package users_repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type UserRepository struct {
    querier *db.Queries
}

func (r *UserRepository) CreateUser(ctx context.Context, username, password string) error {
    id := uuid.New().String()
    return r.querier.CreateUser(ctx, db.CreateUserParams{
        ID:       id,
        Username: username,
        Password: password,
    })
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
    user, err := r.querier.GetUserByUsername(ctx, username)
    if err != nil {
        return domain.User{}, err
    }
    return domain.User{
        ID:       user.ID,
        Username: user.Username,
        Password: user.Password,
    }, nil
}