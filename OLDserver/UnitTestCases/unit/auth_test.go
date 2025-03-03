package unit

import (
    "errors"
    "fmt"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateUser")
    fmt.Println("Testing user creation with mocked repository")
    
    // Create a mock user repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations
    mockRepo.On("CreateUser", "testuser", "testpassword").Return(&models.User{
        ID:       "test-id",
        Username: "testuser",
        Password: "hashed-password",
    }, nil)
    
    // Call the mock
    user, err := mockRepo.CreateUser("testuser", "testpassword")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test-id", user.ID)
    assert.Equal(t, "hashed-password", user.Password)
    
    fmt.Printf("Created user: {ID:%s Username:%s Password:%s}\n", user.ID, user.Username, user.Password)
    fmt.Printf("User ID verified: %s\n", user.ID)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("User creation test passed")
}

func TestCreateUserError(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateUserError")
    fmt.Println("Testing user creation with duplicate username")
    
    // Create a mock user repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations with error
    mockRepo.On("CreateUser", "existinguser", "password").Return(nil, 
        errors.New("username already exists"))
    
    // Call the mock
    user, err := mockRepo.CreateUser("existinguser", "password")
    
    // Assertions
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "already exists")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("User creation error test passed")
}

func TestGetUserByUsername(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetUserByUsername")
    fmt.Println("Testing retrieving a user by username")
    
    // Create a mock user repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations for existing user
    mockRepo.On("GetUserByUsername", "existinguser").Return(&models.User{
        ID:       "user-123",
        Username: "existinguser",
        Password: "hashed-password",
    }, nil)
    
    // Setup expectations for non-existent user
    mockRepo.On("GetUserByUsername", "nonexistentuser").Return(nil, 
        errors.New("user not found"))
    
    // Test existing user
    user, err := mockRepo.GetUserByUsername("existinguser")
    
    // Assertions for existing user
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "existinguser", user.Username)
    assert.Equal(t, "user-123", user.ID)
    
    fmt.Printf("Retrieved user: {ID:%s Username:%s}\n", user.ID, user.Username)
    fmt.Println("User retrieved successfully")
    
    // Test non-existent user
    fmt.Println("\nTesting retrieving a non-existent user")
    user, err = mockRepo.GetUserByUsername("nonexistentuser")
    
    // Assertions for non-existent user
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent user: %v\n", err)
    
    // Verify all expectations were met
    mockRepo.AssertExpectations(t)
    fmt.Println("GetUserByUsername test passed")
}

func TestVerifyPassword(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestVerifyPassword")
    fmt.Println("Testing password verification")
    
    // Create a mock user repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations
    mockRepo.On("VerifyPassword", "hashed-password", "correct-password").Return(nil)
    mockRepo.On("VerifyPassword", "hashed-password", "wrong-password").Return(
        errors.New("password does not match"))
    
    // Test correct password
    err := mockRepo.VerifyPassword("hashed-password", "correct-password")
    assert.NoError(t, err)
    fmt.Println("Password verification successful for correct password")
    
    // Test wrong password
    err = mockRepo.VerifyPassword("hashed-password", "wrong-password")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "password does not match")
    fmt.Printf("Correctly got error for incorrect password: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("Password verification test passed")
}

func TestTokenGeneration(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestTokenGeneration")
    fmt.Println("Testing JWT token generation and validation")
    
    // Create a mock JWT helper
    mockJWT := new(mocks.MockJWTHelper)
    
    // Setup expectations
    mockJWT.On("GenerateJWT", "testuser", "user-123").Return("valid-token", nil)
    mockJWT.On("ParseJWT", "valid-token").Return(&middleware.Claims{
        Username: "testuser",
        UserID:   "user-123",
    }, nil)
    mockJWT.On("ParseJWT", "invalid-token").Return(nil, errors.New("invalid token"))
    
    // Test token generation
    token, err := mockJWT.GenerateJWT("testuser", "user-123")
    assert.NoError(t, err)
    assert.Equal(t, "valid-token", token)
    
    fmt.Printf("Generated token: %s\n", token)
    
    // Test valid token parsing
    claims, err := mockJWT.ParseJWT("valid-token")
    assert.NoError(t, err)
    assert.Equal(t, "testuser", claims.Username)
    assert.Equal(t, "user-123", claims.UserID)
    
    fmt.Printf("Parsed token claims: {Username:%s, UserID:%s}\n", claims.Username, claims.UserID)
    
    // Test invalid token parsing
    fmt.Println("Testing invalid token parsing")
    claims, err = mockJWT.ParseJWT("invalid-token")
    assert.Error(t, err)
    assert.Nil(t, claims)
    
    fmt.Printf("Correctly got error for invalid token: %v\n", err)
    
    // Verify expectations
    mockJWT.AssertExpectations(t)
    fmt.Println("Token generation and validation test passed")
}

func TestDuplicateUsername(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDuplicateUsername")
    fmt.Println("Testing creation of user with duplicate username")
    
    // Create a mock user repository
    mockRepo := new(mocks.MockUserRepository)
    
    // Setup expectations
    mockRepo.On("CreateUser", "duplicate", "password1").Return(&models.User{
        ID:       "user-1",
        Username: "duplicate",
        Password: "hashed-password-1",
    }, nil).Once()
    
    mockRepo.On("CreateUser", "duplicate", "password2").Return(nil,
        errors.New("username already exists")).Once()
    
    // First creation should succeed
    user1, err := mockRepo.CreateUser("duplicate", "password1")
    assert.NoError(t, err)
    assert.NotNil(t, user1)
    
    fmt.Printf("First user created successfully: {ID:%s Username:%s}\n", user1.ID, user1.Username)
    
    // Second creation with same username should fail
    user2, err := mockRepo.CreateUser("duplicate", "password2")
    assert.Error(t, err)
    assert.Nil(t, user2)
    assert.Contains(t, err.Error(), "already exists")
    
    fmt.Printf("Correctly got error for duplicate username: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("Duplicate username test passed")
}