package users_repository

import (
    "context"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

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