package mocks

import (
    "context"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of user service
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.CreateResponse), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
    args := m.Called(ctx, username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) VerifyPassword(hashedPassword, password string) error {
    args := m.Called(hashedPassword, password)
    return args.Error(0)
}

// MockTodoService is a mock implementation of todo service
type MockTodoService struct {
    mock.Mock
}

func (m *MockTodoService) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.CreateResponse), args.Error(1)
}

func (m *MockTodoService) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TodosResponse), args.Error(1)
}

func (m *MockTodoService) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockTodoService) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, id, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockTodoService) UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, id, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockTodoService) GetTodoByID(ctx context.Context, id string) (*domain.Todo, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Todo), args.Error(1)
}

// MockSharedTodoService is a mock implementation of shared todo service
type MockSharedTodoService struct {
    mock.Mock
}

func (m *MockSharedTodoService) GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SharedTodosResponse), args.Error(1)
}

func (m *MockSharedTodoService) GetSharedByMeTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SharedTodosResponse), args.Error(1)
}

func (m *MockSharedTodoService) ShareTodo(ctx context.Context, todoID string, recipientUserID string, sharedBy string) error {
    args := m.Called(ctx, todoID, recipientUserID, sharedBy)
    return args.Error(0)
}

// MockTeamService is a mock implementation of team service
type MockTeamService struct {
    mock.Mock
}

func (m *MockTeamService) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.CreateResponse), args.Error(1)
}

func (m *MockTeamService) GetTeams(ctx context.Context, userID string) (*dto.TeamsResponse, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TeamsResponse), args.Error(1)
}

// MockTeamMemberService is a mock implementation of team member service
type MockTeamMemberService struct {
    mock.Mock
}

func (m *MockTeamMemberService) AddTeamMember(ctx context.Context, req *dto.AddTeamMemberRequest) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockTeamMemberService) GetTeamMembers(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
    args := m.Called(ctx, teamID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TeamMembersResponse), args.Error(1)
}

func (m *MockTeamMemberService) RemoveTeamMember(ctx context.Context, teamID, userID, adminID string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, teamID, userID, adminID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

// MockTeamTodoService is a mock implementation of team todo service
type MockTeamTodoService struct {
    mock.Mock
}

func (m *MockTeamTodoService) CreateTeamTodo(ctx context.Context, req *dto.CreateTeamTodoRequest) (*dto.CreateResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.CreateResponse), args.Error(1)
}

func (m *MockTeamTodoService) GetTeamTodos(ctx context.Context, teamID string) (*dto.TeamTodosResponse, error) {
    args := m.Called(ctx, teamID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TeamTodosResponse), args.Error(1)
}

func (m *MockTeamTodoService) UpdateTeamTodo(ctx context.Context, req *dto.UpdateTeamTodoRequest) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockTeamTodoService) DeleteTeamTodo(ctx context.Context, id, teamID string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, id, teamID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

// MockRoutineService is a mock implementation of routine service
type MockRoutineService struct {
    mock.Mock
}

func (m *MockRoutineService) GetRoutinesByTaskID(ctx context.Context, taskID string) (*dto.RoutinesResponse, error) {
    args := m.Called(ctx, taskID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.RoutinesResponse), args.Error(1)
}

func (m *MockRoutineService) GetDailyRoutines(ctx context.Context, day, scheduleType, userID string) (*dto.TodosResponse, error) {
    args := m.Called(ctx, day, scheduleType, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TodosResponse), args.Error(1)
}

func (m *MockRoutineService) GetTodayRoutines(ctx context.Context, scheduleType, userID string) (*dto.TodosResponse, error) {
    args := m.Called(ctx, scheduleType, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TodosResponse), args.Error(1)
}

func (m *MockRoutineService) UpdateRoutineStatus(ctx context.Context, id string, isActive bool) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, id, isActive)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockRoutineService) UpdateRoutineDay(ctx context.Context, id, day string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, id, day)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}

func (m *MockRoutineService) DeleteRoutinesByTaskID(ctx context.Context, taskID string) (*dto.SuccessResponse, error) {
    args := m.Called(ctx, taskID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.SuccessResponse), args.Error(1)
}