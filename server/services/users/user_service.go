package users

import (
    "context"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type UserService struct {
    repo domain.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error) {
    const functionName = "services.users.UserService.CreateUser"
    
    // Hash the password before storing
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to hash password: %w", functionName, err)
    }
    
    id, err := s.repo.CreateUser(ctx, req.Username, string(hashedPassword))
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create user: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
    const functionName = "services.users.UserService.GetUserByUsername"
    user, err := s.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get user by username: %w", functionName, err)
    }
    
    return &dto.UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Password: user.Password,
    }, nil
}

// VerifyPassword compares a hashed password with a plaintext password
func (s *UserService) VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}