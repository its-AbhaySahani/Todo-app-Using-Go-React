package mocks

import (
    "context"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, password string) (string, error) {
    args := m.Called(ctx, username, password)
    return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
    args := m.Called(ctx, username)
    return args.Get(0).(domain.User), args.Error(1)
}

// MockTodoRepository is a mock implementation of domain.TodoRepository
type MockTodoRepository struct {
    mock.Mock
}

func (m *MockTodoRepository) CreateTodo(ctx context.Context, task, description string, done, important bool, userID string, date, todoTime time.Time) (string, error) {
    args := m.Called(ctx, task, description, done, important, userID, date, todoTime)
    return args.String(0), args.Error(1)
}

func (m *MockTodoRepository) GetTodosByUserID(ctx context.Context, userID string) ([]domain.Todo, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) (bool, error) {
    args := m.Called(ctx, id, task, description, done, important, userID)
    return args.Bool(0), args.Error(1)
}

func (m *MockTodoRepository) DeleteTodo(ctx context.Context, id, userID string) (bool, error) {
    args := m.Called(ctx, id, userID)
    return args.Bool(0), args.Error(1)
}

func (m *MockTodoRepository) UndoTodo(ctx context.Context, id, userID string) (bool, error) {
    args := m.Called(ctx, id, userID)
    return args.Bool(0), args.Error(1)
}

func (m *MockTodoRepository) GetTodoByID(ctx context.Context, id string) (*domain.Todo, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Todo), args.Error(1)
}

// MockSharedTodoRepository is a mock implementation of domain.SharedTodoRepository
type MockSharedTodoRepository struct {
    mock.Mock
}

func (m *MockSharedTodoRepository) CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) (string, error) {
    args := m.Called(ctx, task, description, done, important, userID, sharedBy)
    return args.String(0), args.Error(1)
}

func (m *MockSharedTodoRepository) GetSharedTodos(ctx context.Context, userID string) ([]domain.SharedTodo, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]domain.SharedTodo), args.Error(1)
}

func (m *MockSharedTodoRepository) GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]domain.SharedTodo, error) {
    args := m.Called(ctx, sharedBy)
    return args.Get(0).([]domain.SharedTodo), args.Error(1)
}

func (m *MockSharedTodoRepository) ShareTodo(ctx context.Context, originalTodoID string, recipientUserID string, sharedBy string) error {
    args := m.Called(ctx, originalTodoID, recipientUserID, sharedBy)
    return args.Error(0)
}

func (m *MockSharedTodoRepository) IsSharedWithUser(ctx context.Context, todoID string, userID string) (bool, error) {
    args := m.Called(ctx, todoID, userID)
    return args.Bool(0), args.Error(1)
}

// MockTeamRepository is a mock implementation of domain.TeamRepository
type MockTeamRepository struct {
    mock.Mock
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, name, password, adminID string) (string, error) {
    args := m.Called(ctx, name, password, adminID)
    return args.String(0), args.Error(1)
}

func (m *MockTeamRepository) GetTeamsByAdminID(ctx context.Context, adminID string) ([]domain.Team, error) {
    args := m.Called(ctx, adminID)
    return args.Get(0).([]domain.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeamByID(ctx context.Context, teamID string) (*domain.Team, error) {
    args := m.Called(ctx, teamID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeams(ctx context.Context, userID string) ([]domain.Team, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]domain.Team), args.Error(1)
}

func (m *MockTeamRepository) JoinTeam(ctx context.Context, teamName, password, userID string) error {
    args := m.Called(ctx, teamName, password, userID)
    return args.Error(0)
}

// MockTeamMemberRepository is a mock implementation of domain.TeamMemberRepository
type MockTeamMemberRepository struct {
    mock.Mock
}

func (m *MockTeamMemberRepository) AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) error {
    args := m.Called(ctx, teamID, userID, isAdmin)
    return args.Error(0)
}

func (m *MockTeamMemberRepository) GetTeamMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error) {
    args := m.Called(ctx, teamID)
    return args.Get(0).([]domain.TeamMember), args.Error(1)
}

func (m *MockTeamMemberRepository) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
    args := m.Called(ctx, teamID, userID)
    return args.Error(0)
}

func (m *MockTeamMemberRepository) IsTeamAdmin(ctx context.Context, teamID, userID string) (bool, error) {
    args := m.Called(ctx, teamID, userID)
    return args.Bool(0), args.Error(1)
}

// MockTeamTodoRepository is a mock implementation of domain.TeamTodoRepository
type MockTeamTodoRepository struct {
    mock.Mock
}

func (m *MockTeamTodoRepository) CreateTeamTodo(ctx context.Context, task, description string, done, important bool, teamID, assignedTo string) (string, error) {
    args := m.Called(ctx, task, description, done, important, teamID, assignedTo)
    return args.String(0), args.Error(1)
}

func (m *MockTeamTodoRepository) GetTeamTodos(ctx context.Context, teamID string) ([]domain.TeamTodo, error) {
    args := m.Called(ctx, teamID)
    return args.Get(0).([]domain.TeamTodo), args.Error(1)
}

func (m *MockTeamTodoRepository) UpdateTeamTodo(ctx context.Context, id, task, description string, done, important bool, teamID, assignedTo string) (bool, error) {
    args := m.Called(ctx, id, task, description, done, important, teamID, assignedTo)
    return args.Bool(0), args.Error(1)
}

func (m *MockTeamTodoRepository) DeleteTeamTodo(ctx context.Context, id, teamID string) (bool, error) {
    args := m.Called(ctx, id, teamID)
    return args.Bool(0), args.Error(1)
}

// MockRoutineRepository is a mock implementation of domain.RoutineRepository
type MockRoutineRepository struct {
    mock.Mock
}

func (m *MockRoutineRepository) CreateRoutine(ctx context.Context, day, scheduleType, taskID, userID string, isActive bool) (string, error) {
    args := m.Called(ctx, day, scheduleType, taskID, userID, isActive)
    return args.String(0), args.Error(1)
}

func (m *MockRoutineRepository) GetRoutinesByTaskID(ctx context.Context, taskID string) ([]domain.Routine, error) {
    args := m.Called(ctx, taskID)
    return args.Get(0).([]domain.Routine), args.Error(1)
}

func (m *MockRoutineRepository) GetRoutinesByUserIDAndDay(ctx context.Context, userID, day, scheduleType string) ([]domain.Routine, error) {
    args := m.Called(ctx, userID, day, scheduleType)
    return args.Get(0).([]domain.Routine), args.Error(1)
}

func (m *MockRoutineRepository) UpdateRoutineStatus(ctx context.Context, id string, isActive bool) (bool, error) {
    args := m.Called(ctx, id, isActive)
    return args.Bool(0), args.Error(1)
}

func (m *MockRoutineRepository) UpdateRoutineDay(ctx context.Context, id, day string) (bool, error) {
    args := m.Called(ctx, id, day)
    return args.Bool(0), args.Error(1)
}

func (m *MockRoutineRepository) DeleteRoutinesByTaskID(ctx context.Context, taskID string) (bool, error) {
    args := m.Called(ctx, taskID)
    return args.Bool(0), args.Error(1)
}