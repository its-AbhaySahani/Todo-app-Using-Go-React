package unit

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMiddleware(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAuthMiddleware")
    fmt.Println("Testing authentication middleware")
    
    // Create a mock auth middleware
    mockAuth := new(mocks.MockAuthMiddleware)
    
    // Create a test handler to check if middleware passes request
    testHandlerCalled := false
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check that userID was added to context
        userID, ok := r.Context().Value(middleware.UserIDKey).(string)
        assert.True(t, ok, "userID should be in context")
        assert.Equal(t, "test-user-id", userID)
        testHandlerCalled = true
        w.WriteHeader(http.StatusOK)
    })
    
    // Test successful authentication
    fmt.Println("Testing successful authentication")
    mockAuth.On("AuthMiddleware", mock.Anything).Return().Once()
    mockAuth.Auth = true
    mockAuth.UserID = "test-user-id"
    
    // Apply middleware and execute request
    req, _ := http.NewRequest("GET", "/api/test", nil)
    rr := httptest.NewRecorder()
    handler := mockAuth.AuthMiddleware(testHandler)
    handler.ServeHTTP(rr, req)
    
    // Assertions for successful auth
    assert.True(t, testHandlerCalled, "The test handler should have been called")
    assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")
    
    fmt.Println("Authentication middleware successfully passed the request to the handler")
    
    // Test failed authentication
    fmt.Println("Testing failed authentication")
    mockAuth.On("AuthMiddleware", mock.Anything).Return().Once()
    mockAuth.Auth = false
    
    // Reset for failed auth test
    testHandlerCalled = false
    req2, _ := http.NewRequest("GET", "/api/test", nil)
    rr2 := httptest.NewRecorder()
    
    // Apply middleware and execute request
    handler2 := mockAuth.AuthMiddleware(testHandler)
    handler2.ServeHTTP(rr2, req2)

	// Assertions for failed auth
    assert.False(t, testHandlerCalled, "The test handler should not have been called")
    assert.Equal(t, http.StatusUnauthorized, rr2.Code, "Status code should be 401 Unauthorized")
    
    fmt.Println("Authentication middleware correctly blocked unauthorized request")
    
    // Verify all expectations were met
    mockAuth.AssertExpectations(t)
    fmt.Println("AuthMiddleware test passed")
}

func TestJWTGeneration(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestJWTGeneration")
    fmt.Println("Testing JWT token generation")
    
    // Create a mock JWT helper
    mockJWT := new(mocks.MockJWTHelper)
    
    // Setup expectations
    mockJWT.On("GenerateJWT", "testuser", "user-123").Return("valid-token", nil)
    
    // Generate token
    token, err := mockJWT.GenerateJWT("testuser", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "valid-token", token)
    
    fmt.Printf("Successfully generated JWT token: %s\n", token)
    
    // Verify expectations
    mockJWT.AssertExpectations(t)
    fmt.Println("JWT generation test passed")
}

func TestJWTValidation(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestJWTValidation")
    fmt.Println("Testing JWT token validation")
    
    // Create a mock JWT helper
    mockJWT := new(mocks.MockJWTHelper)
    
    // Setup expectations for valid token
    mockJWT.On("ParseJWT", "valid-token").Return(&middleware.Claims{
        Username: "testuser",
        UserID:   "user-123",
    }, nil)
    
    // Setup expectations for invalid token
    mockJWT.On("ParseJWT", "invalid-token").Return(nil, fmt.Errorf("invalid token"))
    
    // Test valid token
    claims, err := mockJWT.ParseJWT("valid-token")
    
    // Assertions for valid token
    assert.NoError(t, err)
    assert.NotNil(t, claims)
    assert.Equal(t, "testuser", claims.Username)
    assert.Equal(t, "user-123", claims.UserID)
    
    fmt.Printf("Successfully validated token, extracted claims: Username=%s, UserID=%s\n", 
        claims.Username, claims.UserID)
    
    // Test invalid token
    claims, err = mockJWT.ParseJWT("invalid-token")
    
    // Assertions for invalid token
    assert.Error(t, err)
    assert.Nil(t, claims)
    assert.Contains(t, err.Error(), "invalid token")
    
    fmt.Printf("Correctly rejected invalid token: %v\n", err)
    
    // Verify expectations
    mockJWT.AssertExpectations(t)
    fmt.Println("JWT validation test passed")
}

func TestResponseWriterMock(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestResponseWriterMock")
    fmt.Println("Testing HTTP response writer mock")
    
    // Create a mock response writer
    mockWriter := mocks.NewMockResponseWriter()
    
    // Setup expectations
    mockWriter.On("WriteHeader", http.StatusOK).Return()
    mockWriter.On("Write", []byte("Hello, World!")).Return(13, nil)
    
    // Test the response writer
    mockWriter.WriteHeader(http.StatusOK)
    n, err := mockWriter.Write([]byte("Hello, World!"))
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 13, n)
    assert.Equal(t, http.StatusOK, mockWriter.StatusCode)
    assert.Equal(t, []byte("Hello, World!"), mockWriter.WrittenData)
    
    fmt.Printf("Mock response writer successfully captured status code %d and data: %s\n", 
        mockWriter.StatusCode, string(mockWriter.WrittenData))
    
    // Verify expectations
    mockWriter.AssertExpectations(t)
    fmt.Println("Response writer mock test passed")
}