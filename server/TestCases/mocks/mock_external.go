package mocks

import (
    "database/sql"
    "net/http"

    "github.com/dgrijalva/jwt-go"
    "github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of sql.DB
type MockDB struct {
    mock.Mock
}

func (m *MockDB) Ping() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    mockArgs := m.Called(query, args)
    if mockArgs.Get(0) == nil {
        return nil, mockArgs.Error(1)
    }
    return mockArgs.Get(0).(*sql.Rows), mockArgs.Error(1)
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
    mockArgs := m.Called(query, args)
    return mockArgs.Get(0).(*sql.Row)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    mockArgs := m.Called(query, args)
    return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *MockDB) Begin() (*sql.Tx, error) {
    args := m.Called()
    return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockDB) Close() error {
    args := m.Called()
    return args.Error(0)
}

// MockResponseWriter is a mock implementation of http.ResponseWriter
type MockResponseWriter struct {
    mock.Mock
    StatusCode int
    Body       []byte
    Headers    http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
    return &MockResponseWriter{
        Headers: make(http.Header),
    }
}

func (m *MockResponseWriter) Header() http.Header {
    return m.Headers
}

func (m *MockResponseWriter) Write(body []byte) (int, error) {
    m.Body = append(m.Body, body...)
    args := m.Called(body)
    return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
    m.StatusCode = statusCode
    m.Called(statusCode)
}

// MockClaims is a mock JWT claims struct
type MockClaims struct {
    mock.Mock
    jwt.StandardClaims
    Username string
    UserID   string
}

// MockJWTUtil is a mock for JWT utility functions
type MockJWTUtil struct {
    mock.Mock
}

func (m *MockJWTUtil) GenerateJWT(username, userID string) (string, error) {
    args := m.Called(username, userID)
    return args.String(0), args.Error(1)
}

func (m *MockJWTUtil) ParseJWT(tokenString string) (*MockClaims, error) {
    args := m.Called(tokenString)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*MockClaims), args.Error(1)
}

// MockSQLResult is a mock implementation of sql.Result
type MockSQLResult struct {
    mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
    args := m.Called()
    return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
    args := m.Called()
    return args.Get(0).(int64), args.Error(1)
}