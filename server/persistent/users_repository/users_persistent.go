package users_repository

import (
    "context"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

type UserServiceRepository interface {
    CreateUser(ctx context.Context, username, password string) error
    GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}
type UserRepository struct {
    querier *db.Queries
}


func (r *UserRepository) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateUserDomainRequestToPersistentRequest()
    err := r.querier.CreateUser(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
    user, err := r.querier.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }
    return dto.NewUserResponse(&user), nil
}