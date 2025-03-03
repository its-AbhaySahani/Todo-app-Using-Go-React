package mocks

import (
    "database/sql"
    "github.com/stretchr/testify/mock"
)

// MockDB provides a mock implementation of the database
type MockDB struct {
    mock.Mock
}

// QueryRow mocks the database QueryRow method
func (m *MockDB) QueryRow(query string, args ...interface{}) *MockRow {
    callArgs := m.Called(append([]interface{}{query}, args...)...)
    if callArgs.Get(0) == nil {
        return nil
    }
    return callArgs.Get(0).(*MockRow)
}

// Query mocks the database Query method
func (m *MockDB) Query(query string, args ...interface{}) (*MockRows, error) {
    callArgs := m.Called(append([]interface{}{query}, args...)...)
    if callArgs.Get(0) == nil {
        return nil, callArgs.Error(1)
    }
    return callArgs.Get(0).(*MockRows), callArgs.Error(1)
}

// Exec mocks the database Exec method
func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
    callArgs := m.Called(append([]interface{}{query}, args...)...)
    if callArgs.Get(0) == nil {
        return nil, callArgs.Error(1)
    }
    return callArgs.Get(0).(sql.Result), callArgs.Error(1)
}

// Begin mocks the database Begin method
func (m *MockDB) Begin() (*MockTx, error) {
    callArgs := m.Called()
    if callArgs.Get(0) == nil {
        return nil, callArgs.Error(1)
    }
    return callArgs.Get(0).(*MockTx), callArgs.Error(1)
}

// MockTx is a mock for database transactions
type MockTx struct {
    mock.Mock
}

// Commit mocks the transaction Commit method
func (m *MockTx) Commit() error {
    args := m.Called()
    return args.Error(0)
}

// Rollback mocks the transaction Rollback method
func (m *MockTx) Rollback() error {
    args := m.Called()
    return args.Error(0)
}

// MockRow is a mock for database row
type MockRow struct {
    mock.Mock
    values []interface{}
}

// Scan mocks scanning row values into variables
func (m *MockRow) Scan(dest ...interface{}) error {
    args := m.Called(dest)
    
    // If values are set and scan should succeed, copy them to destinations
    if m.values != nil && args.Error(0) == nil {
        for i, val := range m.values {
            if i < len(dest) {
                switch d := dest[i].(type) {
                case *string:
                    *d = val.(string)
                case *int:
                    *d = val.(int)
                case *bool:
                    *d = val.(bool)
                // Add other types as needed
                }
            }
        }
    }
    
    return args.Error(0)
}

// MockRows is a mock for database rows
type MockRows struct {
    mock.Mock
    currentIndex int
    data         [][]interface{}
    columns      []string
}

// SetData sets the mock data to be returned by the rows
func (m *MockRows) SetData(data [][]interface{}, columns []string) {
    m.data = data
    m.columns = columns
    m.currentIndex = -1
}

// Columns returns the column names
func (m *MockRows) Columns() ([]string, error) {
    args := m.Called()
    return m.columns, args.Error(0)
}

// Next mocks moving to the next row
func (m *MockRows) Next() bool {
    args := m.Called()
    m.currentIndex++
    return m.currentIndex < len(m.data) && args.Bool(0)
}

// Scan mocks scanning row values into variables
func (m *MockRows) Scan(dest ...interface{}) error {
    args := m.Called(dest)
    
    // If data is set and scan should succeed, copy values to destinations
    if m.currentIndex >= 0 && m.currentIndex < len(m.data) && args.Error(0) == nil {
        row := m.data[m.currentIndex]
        for i, val := range row {
            if i < len(dest) {
                switch d := dest[i].(type) {
                case *string:
                    *d = val.(string)
                case *int:
                    *d = val.(int)
                case *bool:
                    *d = val.(bool)
                // Add other types as needed
                }
            }
        }
    }
    
    return args.Error(0)
}

// Close mocks closing the rows
func (m *MockRows) Close() error {
    args := m.Called()
    return args.Error(0)
}

// Err mocks checking for errors after iteration
func (m *MockRows) Err() error {
    args := m.Called()
    return args.Error(0)
}

// MockResult is a mock for sql.Result
type MockResult struct {
    mock.Mock
    lastID      int64
    rowsAffected int64
}

// LastInsertId mocks getting the last inserted ID
func (m *MockResult) LastInsertId() (int64, error) {
    args := m.Called()
    return m.lastID, args.Error(0)
}

// RowsAffected mocks getting the number of rows affected
func (m *MockResult) RowsAffected() (int64, error) {
    args := m.Called()
    return m.rowsAffected, args.Error(0)
}

// NewMockResult creates a new MockResult with specified return values
func NewMockResult(lastID, rowsAffected int64) *MockResult {
    result := &MockResult{
        lastID:      lastID,
        rowsAffected: rowsAffected,
    }
    
    result.On("LastInsertId").Return(lastID, nil)
    result.On("RowsAffected").Return(rowsAffected, nil)
    
    return result
}