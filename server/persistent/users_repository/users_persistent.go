package users_repository

import (
    "context"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure UserRepository implements domain.UserRepository
var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
    querier *db.Queries
}

// Implement domain.UserRepository interface methods
func (r *UserRepository) CreateUser(ctx context.Context, username, password string) (string, error) {
    // Use your existing DTO and converter
    req := &dto.CreateUserRequest{
        Username: username,
        Password: password,
    }
    
    params := req.ConvertCreateUserDomainRequestToPersistentRequest()
    err := r.querier.CreateUser(ctx, *params)
    if err != nil {
        return "", err
    }
    return params.ID, nil
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

// Original methods for backward compatibility
func (r *UserRepository) CreateUserWithDTO(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateUserDomainRequestToPersistentRequest()
    err := r.querier.CreateUser(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *UserRepository) GetUserByUsernameWithDTO(ctx context.Context, username string) (*dto.UserResponse, error) {
    user, err := r.querier.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }
    return dto.NewUserResponse(&user), nil
}