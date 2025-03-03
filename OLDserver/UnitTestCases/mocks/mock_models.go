package mocks

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/stretchr/testify/mock"
)

// MockUserRepository provides mock implementations for user operations
type MockUserRepository struct {
    mock.Mock
}

// CreateUser mocks the user creation operation
func (m *MockUserRepository) CreateUser(username, password string) (*models.User, error) {
    args := m.Called(username, password)
    result := args.Get(0)
    if result == nil {
        return nil, args.Error(1)
    }
    return result.(*models.User), args.Error(1)
}

// GetUserByUsername mocks retrieving a user by username
func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
    args := m.Called(username)
    result := args.Get(0)
    if result == nil {
        return nil, args.Error(1)
    }
    return result.(*models.User), args.Error(1)
}

// VerifyPassword mocks password verification
func (m *MockUserRepository) VerifyPassword(hashedPassword, password string) error {
    args := m.Called(hashedPassword, password)
    return args.Error(0)
}

// MockTodoRepository provides mock implementations for todo operations
type MockTodoRepository struct {
    mock.Mock
}

// GetTodos mocks retrieving all todos for a user
func (m *MockTodoRepository) GetTodos(userID string) ([]models.Todo, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.Todo), args.Error(1)
}

// CreateTodo mocks todo creation
func (m *MockTodoRepository) CreateTodo(task, description string, important bool, userID string) (models.Todo, error) {
    args := m.Called(task, description, important, userID)
    return args.Get(0).(models.Todo), args.Error(1)
}

// UpdateTodo mocks todo updating
func (m *MockTodoRepository) UpdateTodo(id, task, description string, done, important bool, userID string) (models.Todo, error) {
    args := m.Called(id, task, description, done, important, userID)
    return args.Get(0).(models.Todo), args.Error(1)
}

// DeleteTodo mocks todo deletion
func (m *MockTodoRepository) DeleteTodo(id, userID string) error {
    args := m.Called(id, userID)
    return args.Error(0)
}

// UndoTodo mocks marking a todo as not done
func (m *MockTodoRepository) UndoTodo(id, userID string) (models.Todo, error) {
    args := m.Called(id, userID)
    return args.Get(0).(models.Todo), args.Error(1)
}

// MockSharedTodoRepository provides mock implementations for shared todo operations
type MockSharedTodoRepository struct {
    mock.Mock
}

// ShareTodoWithUser mocks sharing a todo
func (m *MockSharedTodoRepository) ShareTodoWithUser(taskID, userID, sharedBy string) error {
    args := m.Called(taskID, userID, sharedBy)
    return args.Error(0)
}

// GetSharedTodos mocks retrieving shared todos for a user
func (m *MockSharedTodoRepository) GetSharedTodos(userID string) ([]models.SharedTodo, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.SharedTodo), args.Error(1)
}

// GetSharedByMeTodos mocks retrieving todos shared by a user
func (m *MockSharedTodoRepository) GetSharedByMeTodos(userID string) ([]models.SharedTodo, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.SharedTodo), args.Error(1)
}

// MockTeamRepository provides mock implementations for team operations
type MockTeamRepository struct {
    mock.Mock
}

// CreateTeam mocks team creation
func (m *MockTeamRepository) CreateTeam(name, password, adminID string) (models.Team, error) {
    args := m.Called(name, password, adminID)
    return args.Get(0).(models.Team), args.Error(1)
}

// JoinTeam mocks joining a team
func (m *MockTeamRepository) JoinTeam(teamName, password, userID string) error {
    args := m.Called(teamName, password, userID)
    return args.Error(0)
}

// GetTeams mocks retrieving teams for a user
func (m *MockTeamRepository) GetTeams(userID string) ([]models.Team, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.Team), args.Error(1)
}

// GetTeamByID mocks retrieving a team by ID
func (m *MockTeamRepository) GetTeamByID(teamID string) (models.Team, error) {
    args := m.Called(teamID)
    return args.Get(0).(models.Team), args.Error(1)
}

// MockTeamMemberRepository provides mock implementations for team member operations
type MockTeamMemberRepository struct {
    mock.Mock
}

// AddTeamMember mocks adding a member to a team
func (m *MockTeamMemberRepository) AddTeamMember(teamID, username string, adminID string) error {
    args := m.Called(teamID, username, adminID)
    return args.Error(0)
}

// GetTeamMembers mocks retrieving all members of a team
func (m *MockTeamMemberRepository) GetTeamMembers(teamID string) ([]models.TeamMemberDetails, error) {
    args := m.Called(teamID)
    return args.Get(0).([]models.TeamMemberDetails), args.Error(1)
}

// RemoveTeamMember mocks removing a member from a team
func (m *MockTeamMemberRepository) RemoveTeamMember(teamID, userID string) error {
    args := m.Called(teamID, userID)
    return args.Error(0)
}

// MockTeamTodoRepository provides mock implementations for team todo operations
type MockTeamTodoRepository struct {
    mock.Mock
}

// CreateTeamTodo mocks team todo creation
func (m *MockTeamTodoRepository) CreateTeamTodo(task, description string, important bool, teamID, assignedTo string) (models.TeamTodo, error) {
    args := m.Called(task, description, important, teamID, assignedTo)
    return args.Get(0).(models.TeamTodo), args.Error(1)
}

// GetTeamTodos mocks retrieving team todos
func (m *MockTeamTodoRepository) GetTeamTodos(teamID string) ([]models.TeamTodo, error) {
    args := m.Called(teamID)
    return args.Get(0).([]models.TeamTodo), args.Error(1)
}

// UpdateTeamTodo mocks updating a team todo
func (m *MockTeamTodoRepository) UpdateTeamTodo(id, task, description string, done, important bool, teamID, assignedTo string) (models.TeamTodo, error) {
    args := m.Called(id, task, description, done, important, teamID, assignedTo)
    return args.Get(0).(models.TeamTodo), args.Error(1)
}

// DeleteTeamTodo mocks deleting a team todo
func (m *MockTeamTodoRepository) DeleteTeamTodo(id, teamID string) error {
    args := m.Called(id, teamID)
    return args.Error(0)
}

// MockRoutineRepository provides mock implementations for routine operations
type MockRoutineRepository struct {
    mock.Mock
}

// GetDailyRoutines mocks retrieving routines for a specific day
func (m *MockRoutineRepository) GetDailyRoutines(day, scheduleType, userID string) ([]models.Todo, error) {
    args := m.Called(day, scheduleType, userID)
    return args.Get(0).([]models.Todo), args.Error(1)
}

// UpdateRoutineDay mocks updating a routine's day
func (m *MockRoutineRepository) UpdateRoutineDay(id, day string) error {
    args := m.Called(id, day)
    return args.Error(0)
}

func (m *MockRoutineRepository) UpdateRoutineStatus(id string, isActive bool) error {
    args := m.Called(id, isActive)
    return args.Error(0)
}

// DeleteRoutinesByTaskID mocks deleting all routines for a task
func (m *MockRoutineRepository) DeleteRoutinesByTaskID(taskID string) error {
    args := m.Called(taskID)
    return args.Error(0)
}

// CreateRoutine mocks creating a routine
func (m *MockRoutineRepository) CreateRoutine(day, scheduleType, taskID, userID string) (models.Routine, error) {
    args := m.Called(day, scheduleType, taskID, userID)
    if args.Get(0) == nil {
        return models.Routine{}, args.Error(1)
    }
    return args.Get(0).(models.Routine), args.Error(1)
}