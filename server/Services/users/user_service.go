package users

import (
    "context"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/users_repository"
)

type UserService struct {
    repo *users_repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error) {
    const functionName = "services.users.UserService.CreateUser"
    res, err := s.repo.CreateUser(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create user: %w", functionName, err)
    }
    return res, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
    const functionName = "services.users.UserService.GetUserByUsername"
    res, err := s.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get user by username: %w", functionName, err)
    }
    return res, nil
}

func (s *UserService) VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}