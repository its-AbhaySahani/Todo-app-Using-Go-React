package mocks

import (
    "context"
    "net/http"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/stretchr/testify/mock"
)

// MockAuthMiddleware provides a mock implementation of the auth middleware
type MockAuthMiddleware struct {
    mock.Mock
    Auth bool
    UserID string
}

// AuthMiddleware mocks the authentication middleware
func (m *MockAuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
    // Record the method call with the handler
    m.Called(next)
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // This function is executed when ServeHTTP is called on the returned handler
        
        // Check if authentication should succeed based on our test expectations
        if !m.Auth {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        // Add user ID to context and call next handler
        ctx := context.WithValue(r.Context(), middleware.UserIDKey, m.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// MockClaims provides a mock for JWT claims
type MockClaims struct {
    Username string
    UserID   string
    mock.Mock
}

// MockJWTHelper provides mock implementations of JWT-related functions
type MockJWTHelper struct {
    mock.Mock
}

// GenerateJWT mocks the JWT generation process
func (m *MockJWTHelper) GenerateJWT(username, userID string) (string, error) {
    args := m.Called(username, userID)
    return args.String(0), args.Error(1)
}

// ParseJWT mocks JWT token parsing
func (m *MockJWTHelper) ParseJWT(tokenString string) (*middleware.Claims, error) {
    args := m.Called(tokenString)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*middleware.Claims), args.Error(1)
}

// MockResponseWriter is a mock implementation of http.ResponseWriter
type MockResponseWriter struct {
    mock.Mock
    Headers     http.Header
    WrittenData []byte
    StatusCode  int
}

// NewMockResponseWriter creates a new MockResponseWriter with default values
func NewMockResponseWriter() *MockResponseWriter {
    return &MockResponseWriter{
        Headers:    make(http.Header),
        StatusCode: http.StatusOK,
    }
}

func (m *MockResponseWriter) Header() http.Header {
    return m.Headers
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
    args := m.Called(data)
    m.WrittenData = append(m.WrittenData, data...)
    return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
    m.Called(statusCode)
    m.StatusCode = statusCode
}