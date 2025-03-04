package services_test

import (
    "context"
    "errors"
    "fmt"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/users"
    "github.com/stretchr/testify/assert"
    "golang.org/x/crypto/bcrypt"
    "github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestCreateUser ===")
    fmt.Println("Testing user creation functionality")
    
    // Create a mock repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repository
    userService := users.NewUserService(mockRepo)
    
    // Test successful user creation
    mockRepo.On("CreateUser", context.Background(), "testuser", mock.AnythingOfType("string")).Return("user-123", nil)
    
    req := &dto.CreateUserRequest{
        Username: "testuser",
        Password: "password123",
    }
    
    fmt.Println("Scenario 1: Testing successful user creation")
    // Test the service
    res, err := userService.CreateUser(context.Background(), req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, "user-123", res.ID)
    fmt.Printf("✅ User created successfully with ID: %s\n", res.ID)
    
    // Test user creation failure
    mockRepo.On("CreateUser", context.Background(), "existinguser", mock.AnythingOfType("string")).
        Return("", errors.New("username already exists"))
    
    req = &dto.CreateUserRequest{
        Username: "existinguser",
        Password: "password456",
    }
    
    fmt.Println("\nScenario 2: Testing user creation with duplicate username")
    // Test the service
    res, err = userService.CreateUser(context.Background(), req)
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Contains(t, err.Error(), "username already exists")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify that all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All CreateUser test scenarios passed")
}

func TestGetUserByUsername(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetUserByUsername ===")
    fmt.Println("Testing retrieving user by username")
    
    // Create a mock repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repository
    userService := users.NewUserService(mockRepo)
    
    // Set up the mock repository's expectations
    mockUser := domain.User{
        ID:       "user-123",
        Username: "testuser",
        Password: "hashed:password123",
    }
    
    mockRepo.On("GetUserByUsername", context.Background(), "testuser").Return(mockUser, nil)
    mockRepo.On("GetUserByUsername", context.Background(), "nonexistent").Return(domain.User{}, errors.New("user not found"))
    
    // Test getting an existing user
    fmt.Println("Scenario 1: Testing retrieving an existing user")
    user, err := userService.GetUserByUsername(context.Background(), "testuser")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "user-123", user.ID)
    assert.Equal(t, "testuser", user.Username)
    fmt.Printf("✅ Successfully retrieved user: {ID: %s, Username: %s}\n", user.ID, user.Username)
    
    // Test getting a non-existent user
    fmt.Println("\nScenario 2: Testing retrieving a non-existent user")
    user, err = userService.GetUserByUsername(context.Background(), "nonexistent")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "user not found")
    fmt.Printf("✅ Correctly received error: %v\n", err)
    
    // Verify that all expected methods were called
    mockRepo.AssertExpectations(t)
    fmt.Println("✅ All GetUserByUsername test scenarios passed")
}

func TestVerifyPassword(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestVerifyPassword ===")
    fmt.Println("Testing password verification functionality")
    
    mockRepo := new(mocks.MockUserRepository)
    
    // Create the service with the mock repository
    userService := users.NewUserService(mockRepo)
    
    // Generate a real hash for testing
    password := "password123"
    validHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    assert.NoError(t, err, "Failed to generate bcrypt hash")
    
    // Option 1: Test with a freshly generated valid hash for the password
    fmt.Println("Scenario 1: Testing with correct password")
    err = userService.VerifyPassword(string(validHash), password)
    // Hash should be correct for the password
    assert.NoError(t, err)
    fmt.Println("✅ Password verification successful")
    
    // Option 2: Test with a valid hash but wrong password
    fmt.Println("\nScenario 2: Testing with incorrect password")
    err = userService.VerifyPassword(string(validHash), "wrongpassword")
    // Should get an error since password doesn't match hash
    assert.Error(t, err)
    fmt.Printf("✅ Correctly received error for incorrect password: %v\n", err)
    
    // Option 3: Verify the method handles invalid hash formats appropriately
    fmt.Println("\nScenario 3: Testing with invalid hash format")
    invalidHash := "not-a-valid-hash"
    err = userService.VerifyPassword(invalidHash, password)
    // Should get an error about invalid hash format
    assert.Error(t, err)
    fmt.Printf("✅ Correctly received error for invalid hash: %v\n", err)
    
    fmt.Println("✅ All VerifyPassword test scenarios passed")
}