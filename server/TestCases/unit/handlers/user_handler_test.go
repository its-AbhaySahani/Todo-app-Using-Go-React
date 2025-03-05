package handlers_test

import (
    "context"
    "errors"
    "fmt"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// UserServiceAdapterRepo adapts MockUserService to satisfy domain.UserRepository
type UserServiceAdapterRepo struct {
    Mock *mocks.MockUserService
}

// Define the handler logic functions to be tested

// registerHandlerLogic processes user registration logic without HTTP concerns
func registerHandlerLogic(ctx context.Context, username, password string, service userServicer) (*dto.CreateResponse, error) {
    // Validate input
    if username == "" || password == "" {
        return nil, errors.New("username and password are required")
    }
    
    // Create user request
    req := &dto.CreateUserRequest{
        Username: username,
        Password: password,
    }
    
    // Call service
    return service.CreateUser(ctx, req)
}

// loginHandlerLogic processes user login logic without HTTP concerns
func loginHandlerLogic(ctx context.Context, username, password string, service userServicer) (string, string, error) {
    // Validate input
    if username == "" || password == "" {
        return "", "", errors.New("username and password are required")
    }
    
    // Get user by username
    user, err := service.GetUserByUsername(ctx, username)
    if err != nil {
        return "", "", err
    }
    
    // Verify password
    err = service.VerifyPassword(user.Password, password)
    if err != nil {
        return "", "", err
    }
    
    // In an actual implementation, we would generate a JWT token here
    // For testing purposes, we'll return a placeholder token
    token := "valid-jwt-token"
    
    return token, user.ID, nil
}

// getUserProfileLogic processes user profile retrieval logic without HTTP concerns
func getUserProfileLogic(ctx context.Context, userID string, service userServicer) (*domain.User, error) {
    // Get user from service
    return service.GetUserByUsername(ctx, userID)
}

// Interface for user service to allow mocking
type userServicer interface {
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error)
    GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
    VerifyPassword(hashedPassword, password string) error
}

func TestRegisterHandlerLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestRegisterHandlerLogic ===")
    fmt.Println("Testing user registration logic")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    // Setup test data
    username := "testuser"
    password := "password123"
    userId := "user-123"
    
    // Setup expectations - successful registration
    mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *dto.CreateUserRequest) bool {
        return req.Username == username // Don't check password as it might be hashed
    })).Return(&dto.CreateResponse{ID: userId}, nil)
    
    // Test successful registration
    fmt.Println("Scenario 1: Testing successful user registration")
    res, err := registerHandlerLogic(context.Background(), username, password, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, userId, res.ID)
    fmt.Printf("✅ User registered successfully with ID: %s\n", res.ID)
    
    // Setup expectations - existing username
    mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *dto.CreateUserRequest) bool {
        return req.Username == "existinguser"
    })).Return(nil, errors.New("username already exists"))
    
    // Test duplicate username
    fmt.Println("\nScenario 2: Testing registration with existing username")
    res, err = registerHandlerLogic(context.Background(), "existinguser", "password123", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "username already exists")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty username
    fmt.Println("\nScenario 3: Testing registration with empty username")
    res, err = registerHandlerLogic(context.Background(), "", "password123", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "username and password are required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All RegisterHandlerLogic test scenarios passed")
}

func TestLoginHandlerLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestLoginHandlerLogic ===")
    fmt.Println("Testing user login logic")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    // Setup test data
    username := "testuser"
    password := "password123"
    userId := "user-123"
    
    // Setup mock user
    mockUser := &domain.User{
        ID:       userId,
        Username: username,
        Password: "hashed_password", // Would be a bcrypt hash in reality
    }
    
    // Setup expectations - successful login
    mockService.On("GetUserByUsername", mock.Anything, username).Return(mockUser, nil)
    mockService.On("VerifyPassword", "hashed_password", password).Return(nil)
    
    // Test successful login
    fmt.Println("Scenario 1: Testing successful login")
    token, id, err := loginHandlerLogic(context.Background(), username, password, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.Equal(t, userId, id)
    fmt.Printf("✅ User logged in successfully, got token and ID: %s\n", id)
    
    // Setup expectations - non-existent user
    mockService.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("user not found"))
    
    // Test invalid username
    fmt.Println("\nScenario 2: Testing login with invalid username")
    token, id, err = loginHandlerLogic(context.Background(), "nonexistent", password, mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Empty(t, token)
    assert.Empty(t, id)
    assert.Contains(t, err.Error(), "user not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Setup expectations - wrong password
    mockService.On("GetUserByUsername", mock.Anything, "validuser").Return(mockUser, nil)
    mockService.On("VerifyPassword", "hashed_password", "wrongpass").Return(errors.New("invalid password"))
    
    // Test wrong password
    fmt.Println("\nScenario 3: Testing login with wrong password")
    token, id, err = loginHandlerLogic(context.Background(), "validuser", "wrongpass", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Empty(t, token)
    assert.Empty(t, id)
    assert.Contains(t, err.Error(), "invalid password")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty credentials
    fmt.Println("\nScenario 4: Testing login with empty credentials")
    token, id, err = loginHandlerLogic(context.Background(), "", "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Empty(t, token)
    assert.Empty(t, id)
    assert.Contains(t, err.Error(), "username and password are required")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All LoginHandlerLogic test scenarios passed")
}

func TestGetUserProfileLogic(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetUserProfileLogic ===")
    fmt.Println("Testing user profile retrieval logic")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    // Setup test data
    username := "testuser"
    userId := "user-123"
    
    // Setup mock user
    mockUser := &domain.User{
        ID:       userId,
        Username: username,
        Password: "hashed_password", // Would be a bcrypt hash in reality
    }
    
    // Setup expectations - successful profile retrieval
    mockService.On("GetUserByUsername", mock.Anything, userId).Return(mockUser, nil)
    
    // Test successful profile retrieval
    fmt.Println("Scenario 1: Testing successful profile retrieval")
    user, err := getUserProfileLogic(context.Background(), userId, mockService)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, userId, user.ID)
    assert.Equal(t, username, user.Username)
    fmt.Printf("✅ Retrieved user profile for: %s\n", username)
    
    // Setup expectations - non-existent user
    mockService.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("user not found"))
    
    // Test non-existent user
    fmt.Println("\nScenario 2: Testing profile retrieval for non-existent user")
    user, err = getUserProfileLogic(context.Background(), "nonexistent", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "user not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Test empty user ID
    fmt.Println("\nScenario 3: Testing profile retrieval with empty user ID")
    mockService.On("GetUserByUsername", mock.Anything, "").Return(nil, errors.New("invalid user ID"))
    
    user, err = getUserProfileLogic(context.Background(), "", mockService)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "invalid user ID")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All GetUserProfileLogic test scenarios passed")
}