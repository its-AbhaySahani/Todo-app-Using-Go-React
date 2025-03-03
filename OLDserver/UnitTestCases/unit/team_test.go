package unit

import (
    "errors"
    "fmt"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTeam")
    fmt.Println("Testing creating a new team")
    
    // Create a mock team repository
    mockRepo := new(mocks.MockTeamRepository)
    
    // Setup expectations
    mockRepo.On("CreateTeam", "Test Team", "hashedpassword", "admin-123").Return(
        models.Team{
            ID:       "team-1",
            Name:     "Test Team",
            Password: "hashedpassword",
            AdminID:  "admin-123",
        }, nil)
    
    mockRepo.On("CreateTeam", "Existing Team", "hashedpassword", "admin-123").Return(
        models.Team{}, errors.New("team name already exists"))
    
    // Test successful team creation
    team, err := mockRepo.CreateTeam("Test Team", "hashedpassword", "admin-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "team-1", team.ID)
    assert.Equal(t, "Test Team", team.Name)
    assert.Equal(t, "admin-123", team.AdminID)
    
    fmt.Printf("Created team: {ID:%s Name:%s AdminID:%s}\n", team.ID, team.Name, team.AdminID)
    
    // Test team creation with duplicate name
    _, err = mockRepo.CreateTeam("Existing Team", "hashedpassword", "admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already exists")
    
    fmt.Printf("Correctly got error for duplicate team name: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("CreateTeam test passed")
}

func TestJoinTeam(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestJoinTeam")
    fmt.Println("Testing joining an existing team")
    
    // Create a mock team repository
    mockRepo := new(mocks.MockTeamRepository)
    
    // Setup expectations
    mockRepo.On("JoinTeam", "Existing Team", "correctpassword", "user-123").Return(nil)
    mockRepo.On("JoinTeam", "Nonexistent Team", "password", "user-123").Return(
        errors.New("team not found"))
    mockRepo.On("JoinTeam", "Existing Team", "wrongpassword", "user-123").Return(
        errors.New("incorrect password"))
    mockRepo.On("JoinTeam", "Existing Team", "correctpassword", "already-member").Return(
        errors.New("user is already a member of this team"))
    
    // Test successful team join
    err := mockRepo.JoinTeam("Existing Team", "correctpassword", "user-123")
    assert.NoError(t, err)
    
    fmt.Println("Successfully joined team")
    
    // Test joining non-existent team
    err = mockRepo.JoinTeam("Nonexistent Team", "password", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent team: %v\n", err)
    
    // Test joining with incorrect password
    err = mockRepo.JoinTeam("Existing Team", "wrongpassword", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "incorrect password")
    
    fmt.Printf("Correctly got error for incorrect password: %v\n", err)
    
    // Test joining team user is already a member of
    err = mockRepo.JoinTeam("Existing Team", "correctpassword", "already-member")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already a member")
    
    fmt.Printf("Correctly got error for already a member: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("JoinTeam test passed")
}

func TestGetTeams(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeams")
    fmt.Println("Testing retrieving teams for a user")
    
    // Create a mock team repository
    mockRepo := new(mocks.MockTeamRepository)
    
    // Create test teams
    teams := []models.Team{
        {
            ID:       "team-1",
            Name:     "Team One",
            Password: "hashedpassword",
            AdminID:  "admin-123",
        },
        {
            ID:       "team-2",
            Name:     "Team Two",
            Password: "hashedpassword",
            AdminID:  "user-123",
        },
    }
    
    // Setup expectations
    mockRepo.On("GetTeams", "user-123").Return(teams, nil)
    mockRepo.On("GetTeams", "user-no-teams").Return([]models.Team{}, nil)
    mockRepo.On("GetTeams", "error-user").Return([]models.Team{}, errors.New("database error"))
    
    // Test getting teams for user with teams
    userTeams, err := mockRepo.GetTeams("user-123")
    assert.NoError(t, err)
    assert.Equal(t, 2, len(userTeams))
    assert.Equal(t, "Team One", userTeams[0].Name)
    assert.Equal(t, "Team Two", userTeams[1].Name)
    
    fmt.Printf("Retrieved %d teams\n", len(userTeams))
    for i, team := range userTeams {
        fmt.Printf("Team %d: {ID:%s Name:%s AdminID:%s}\n", 
            i+1, team.ID, team.Name, team.AdminID)
    }
    
    // Test getting teams for user with no teams
    emptyTeams, err := mockRepo.GetTeams("user-no-teams")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyTeams))
    
    fmt.Println("Successfully retrieved empty teams list for user with no teams")
    
    // Test database error
    _, err = mockRepo.GetTeams("error-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTeams test passed")
}

func TestGetTeamByID(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamByID")
    fmt.Println("Testing retrieving a team by ID")
    
    // Create a mock team repository
    mockRepo := new(mocks.MockTeamRepository)
    
    // Setup expectations
    mockRepo.On("GetTeamByID", "team-1").Return(
        models.Team{
            ID:       "team-1",
            Name:     "Test Team",
            Password: "hashedpassword",
            AdminID:  "admin-123",
        }, nil)
    
    mockRepo.On("GetTeamByID", "nonexistent-team").Return(
        models.Team{}, errors.New("team not found"))
    
    // Test getting existing team
    team, err := mockRepo.GetTeamByID("team-1")
    assert.NoError(t, err)
    assert.Equal(t, "team-1", team.ID)
    assert.Equal(t, "Test Team", team.Name)
    assert.Equal(t, "admin-123", team.AdminID)
    
    fmt.Printf("Retrieved team: {ID:%s Name:%s AdminID:%s}\n", team.ID, team.Name, team.AdminID)
    
    // Test getting non-existent team
    _, err = mockRepo.GetTeamByID("nonexistent-team")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent team: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTeamByID test passed")
}