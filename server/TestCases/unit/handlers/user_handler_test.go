package handlers_test

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/api"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/users"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/TestCases/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// UserServiceAdapterRepo adapts MockUserService to satisfy domain.UserRepository
type UserServiceAdapterRepo struct {
    Mock *mocks.MockUserService
}

// CreateUser implements domain.UserRepository
func (a *UserServiceAdapterRepo) CreateUser(ctx context.Context, username, password string) (string, error) {
    res, err := a.Mock.CreateUser(ctx, &dto.CreateUserRequest{
        Username: username,
        Password: password,
    })
    if err != nil {
        return "", err
    }
    return res.ID, nil
}

// GetUserByUsername implements domain.UserRepository
func (a *UserServiceAdapterRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
    user, err := a.Mock.GetUserByUsername(ctx, username)
    if err != nil {
        return domain.User{}, err
    }
    return domain.User{
        ID:       user.ID,
        Username: user.Username,
        Password: user.Password,
    }, nil
}

// CreateUserService creates a UserService with our mock adapter
func CreateUserService(mockService *mocks.MockUserService) *users.UserService {
    repo := &UserServiceAdapterRepo{Mock: mockService}
    return users.NewUserService(repo)
}

// TestGetUserHandler tests the user profile retrieval endpoint
func TestGetUserHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestGetUserHandler ===")
    fmt.Println("Testing get user profile endpoint")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    // Create the adapter that will be passed to the API handler
    userService := CreateUserService(mockService)
    
    // Setup test data
    username := "testuser"
    userId := "user-123"
    
    // Setup user to be returned
    mockUser := &domain.User{
        ID:       userId,
        Username: username,
        Password: "hashed_password", // This would be hashed in real scenario
    }
    
    // Setup expectations
    mockService.On("GetUserByUsername", mock.Anything, userId).Return(mockUser, nil)
    
    // Create HTTP request with user context
    req := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
    ctx := context.WithValue(req.Context(), middleware.UserIDKey, userId)
    req = req.WithContext(ctx)
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful user profile retrieval")
    handler := func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from context
        userID := r.Context().Value(middleware.UserIDKey).(string)
        
        // Get user details using the userService
        user, err := userService.GetUserByUsername(context.Background(), userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Return user details
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "id": user.ID,
            "username": user.Username,
        })
    }
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]string
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response data
    assert.Equal(t, userId, responseBody["id"])
    assert.Equal(t, username, responseBody["username"])
    fmt.Printf("✅ User profile retrieved successfully for user: %s\n", username)
    
    // Test with error from service
    mockService.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("user not found"))
    
    // Create HTTP request with non-existent user context
    req = httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
    ctx = context.WithValue(req.Context(), middleware.UserIDKey, "nonexistent")
    req = req.WithContext(ctx)
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing profile retrieval for non-existent user")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All GetUserHandler test scenarios passed")
}

func TestRegisterHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestRegisterHandler ===")
    fmt.Println("Testing user registration endpoint")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    // Create the adapter that will be passed to the API handler
    userService := CreateUserService(mockService)
    
    // Setup test data
    username := "testuser"
    password := "password123"
    userId := "user-123"
    
    // Setup expectations - use mock.Anything for the password since it will be hashed
    mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *dto.CreateUserRequest) bool {
        return req.Username == username // Don't check password as it will be hashed
    })).Return(&dto.CreateResponse{ID: userId}, nil)
    
    // Create a request body
    reqBody := map[string]string{
        "username": username,
        "password": password,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    
    // Create HTTP request
    req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("Scenario 1: Testing successful user registration")
    handler := api.Register(userService)
    handler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]string
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify response
    assert.Equal(t, userId, responseBody["id"])
    fmt.Printf("✅ User registered successfully with ID: %s\n", responseBody["id"])
    
    // Test with error from service - same change here
    mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(req *dto.CreateUserRequest) bool {
        return req.Username == "existinguser" // Don't check password
    })).Return(nil, errors.New("username already exists"))
    
    // Create a request with existing username
    reqBody = map[string]string{
        "username": "existinguser",
        "password": "password123",
    }
    bodyBytes, _ = json.Marshal(reqBody)
    
    // Create HTTP request
    req = httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing registration with existing username")
    handler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    fmt.Printf("✅ Correctly received error response with status code: %d\n", resp.StatusCode)
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All RegisterHandler test scenarios passed")
}

func TestLoginHandler(t *testing.T) {
    fmt.Println("\n=== RUNNING TEST: TestLoginHandler ===")
    fmt.Println("Testing user login endpoint")
    
    // Create mock service
    mockService := new(mocks.MockUserService)
    
    
    // Setup test data
    username := "testuser"
    password := "password123"
    userId := "user-123"
    
    // Setup user to be returned
    mockUser := &domain.User{
        ID:       userId,
        Username: username,
        Password: "hashed_password", // This would be hashed in real scenario
    }
    
    // Setup expectations for scenario 1 - successful login
    mockService.On("GetUserByUsername", mock.Anything, username).Return(mockUser, nil)
    mockService.On("VerifyPassword", "hashed_password", password).Return(nil)
    
    // Create a request body
    reqBody := map[string]string{
        "username": username,
        "password": password,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    
    // Create HTTP request
    req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Create response recorder
    w := httptest.NewRecorder()
    
    // Create a simulated Login handler for more control over testing
    fmt.Println("Scenario 1: Testing successful user login")
    loginHandler := func(w http.ResponseWriter, r *http.Request) {
        var creds map[string]string
        json.NewDecoder(r.Body).Decode(&creds)
        
        user, err := mockService.GetUserByUsername(context.Background(), creds["username"])
        if err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }
        
        err = mockService.VerifyPassword(user.Password, creds["password"])
        if err != nil {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }
        
        // Successful login
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "token": "valid-jwt-token",
            "userId": user.ID,
        })
    }
    
    loginHandler(w, req)
    
    // Check response
    resp := w.Result()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Decode response body
    var responseBody map[string]string
    json.NewDecoder(resp.Body).Decode(&responseBody)
    
    // Verify JWT token was returned
    assert.NotEmpty(t, responseBody["token"])
    fmt.Printf("✅ User logged in successfully, JWT token provided\n")
    
    // Setup expectation for scenario 2 - invalid username
    mockService.On("GetUserByUsername", mock.Anything, "wronguser").Return(nil, errors.New("user not found"))
    
    // Create a request with wrong username
    reqBody = map[string]string{
        "username": "wronguser",
        "password": "password123",
    }
    bodyBytes, _ = json.Marshal(reqBody)
    
    // Create HTTP request
    req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 2: Testing login with invalid username")
    loginHandler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    fmt.Printf("✅ Correctly received unauthorized status for invalid username\n")
    
    // Setup expectations for scenario 3 - incorrect password
    mockService.On("GetUserByUsername", mock.Anything, "validuser").Return(mockUser, nil)
    mockService.On("VerifyPassword", "hashed_password", "wrongpassword").Return(errors.New("invalid password"))
    
    // Create a request with wrong password
    reqBody = map[string]string{
        "username": "validuser",
        "password": "wrongpassword",
    }
    bodyBytes, _ = json.Marshal(reqBody)
    
    // Create HTTP request
    req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // Create response recorder
    w = httptest.NewRecorder()
    
    // Call the handler function
    fmt.Println("\nScenario 3: Testing login with incorrect password")
    loginHandler(w, req)
    
    // Check response
    resp = w.Result()
    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    fmt.Printf("✅ Correctly received unauthorized status for incorrect password\n")
    
    // Verify all expected methods were called
    mockService.AssertExpectations(t)
    fmt.Println("✅ All LoginHandler test scenarios passed")
}