package unit

import (
    "errors"
    "fmt"
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
)

func TestAddTeamMember(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAddTeamMember")
    fmt.Println("Testing adding a member to a team")
    
    // Create a mock team member repository
    mockRepo := new(mocks.MockTeamMemberRepository)
    
    // Setup expectations
    mockRepo.On("AddTeamMember", "team-1", "testuser", "admin-123").Return(nil)
    mockRepo.On("AddTeamMember", "team-1", "nonexistent-user", "admin-123").Return(
        errors.New("user not found"))
    mockRepo.On("AddTeamMember", "nonexistent-team", "testuser", "admin-123").Return(
        errors.New("team not found"))
    mockRepo.On("AddTeamMember", "team-1", "testuser", "non-admin-123").Return(
        errors.New("only admin can add members"))
    mockRepo.On("AddTeamMember", "team-1", "already-member", "admin-123").Return(
        errors.New("user is already a member of this team"))
    
    // Test successful member addition
    err := mockRepo.AddTeamMember("team-1", "testuser", "admin-123")
    assert.NoError(t, err)
    
    fmt.Println("Successfully added team member")
    
    // Test adding non-existent user
    err = mockRepo.AddTeamMember("team-1", "nonexistent-user", "admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "user not found")
    
    fmt.Printf("Correctly got error for non-existent user: %v\n", err)
    
    // Test adding to non-existent team
    err = mockRepo.AddTeamMember("nonexistent-team", "testuser", "admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for non-existent team: %v\n", err)
    
    // Test adding by non-admin
    err = mockRepo.AddTeamMember("team-1", "testuser", "non-admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "only admin can add members")
    
    fmt.Printf("Correctly got error for non-admin: %v\n", err)
    
    // Test adding user that's already a member
    err = mockRepo.AddTeamMember("team-1", "already-member", "admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already a member")
    
    fmt.Printf("Correctly got error for already a member: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("AddTeamMember test passed")
}

func TestGetTeamMembers(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamMembers")
    fmt.Println("Testing retrieving members of a team")
    
    // Create a mock team member repository
    mockRepo := new(mocks.MockTeamMemberRepository)
    
    // Create test members
    members := []models.TeamMemberDetails{
        {
            ID:       "user-1",         // Changed from UserID to ID
            Username: "testuser1",
            IsAdmin:  true,
        },
        {
            ID:       "user-2",         // Changed from UserID to ID
            Username: "testuser2",
            IsAdmin:  false,
        },
    }
    
    // Setup expectations
    mockRepo.On("GetTeamMembers", "team-1").Return(members, nil)
    mockRepo.On("GetTeamMembers", "empty-team").Return([]models.TeamMemberDetails{}, nil)
    mockRepo.On("GetTeamMembers", "nonexistent-team").Return(
        []models.TeamMemberDetails{}, errors.New("team not found"))
    
    // Test getting members of a team with members
    teamMembers, err := mockRepo.GetTeamMembers("team-1")
    assert.NoError(t, err)
    assert.Equal(t, 2, len(teamMembers))
    assert.Equal(t, "testuser1", teamMembers[0].Username)
    assert.Equal(t, "testuser2", teamMembers[1].Username)
    assert.True(t, teamMembers[0].IsAdmin)
    assert.False(t, teamMembers[1].IsAdmin)
    
    fmt.Printf("Retrieved %d team members\n", len(teamMembers))
    for i, member := range teamMembers {
        fmt.Printf("Member %d: {ID:%s Username:%s IsAdmin:%t}\n", 
            i+1, member.ID, member.Username, member.IsAdmin)  // Changed from UserID to ID
    }
    
    // Test getting members of a team with no members
    emptyMembers, err := mockRepo.GetTeamMembers("empty-team")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyMembers))
    
    fmt.Println("Successfully retrieved empty members list for team with no members")
    
    // Test getting members of a non-existent team
    _, err = mockRepo.GetTeamMembers("nonexistent-team")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for non-existent team: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTeamMembers test passed")
}

func TestRemoveTeamMember(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestRemoveTeamMember")
    fmt.Println("Testing removing a member from a team")
    
    // Create a mock team member repository
    mockRepo := new(mocks.MockTeamMemberRepository)
    
    // Setup expectations
    mockRepo.On("RemoveTeamMember", "team-1", "user-123").Return(nil)
    mockRepo.On("RemoveTeamMember", "team-1", "admin-123").Return(
        errors.New("cannot remove the team admin"))
    mockRepo.On("RemoveTeamMember", "nonexistent-team", "user-123").Return(
        errors.New("team not found"))
    mockRepo.On("RemoveTeamMember", "team-1", "nonexistent-user").Return(
        errors.New("user is not a member of this team"))
    
    // Test successful member removal
    err := mockRepo.RemoveTeamMember("team-1", "user-123")
    assert.NoError(t, err)
    
    fmt.Println("Successfully removed team member")
    
    // Test removing team admin
    err = mockRepo.RemoveTeamMember("team-1", "admin-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot remove the team admin")
    
    fmt.Printf("Correctly got error for removing team admin: %v\n", err)
    
    // Test removing from non-existent team
    err = mockRepo.RemoveTeamMember("nonexistent-team", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for non-existent team: %v\n", err)
    
    // Test removing non-member
    err = mockRepo.RemoveTeamMember("team-1", "nonexistent-user")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not a member")
    
    fmt.Printf("Correctly got error for non-member: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("RemoveTeamMember test passed")
}