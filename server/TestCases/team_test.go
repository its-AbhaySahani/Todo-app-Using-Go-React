// server/TestCases/team_test.go

package TestCases

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Note: We're reusing the router and helper functions from team_member_test.go

func TestCreateTeam(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTeam")
    
    // Setup test user
    userID := createTestUser(t)
    
    // Create team request
    team := models.Team{
        Name:        "Test Team Creation",
        Password:    "testpassword",
        AdminID:     userID,
    }
    
    jsonTeam, _ := json.Marshal(team)
    req, _ := http.NewRequest("POST", "/api/team", bytes.NewBuffer(jsonTeam))
    req.Header.Set("Content-Type", "application/json")
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200", w.Code)
    }
    
    // Parse response to get team ID
    var response struct {
        ID string `json:"id"`
    }
    
    t.Logf("Sending request to create team")
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    if response.ID == "" {
        t.Fatalf("Expected a team ID but got empty string")
    }
    
    t.Logf("Created team with ID: %s", response.ID)
    t.Log("TestCreateTeam passed")
}

func TestGetTeams(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeams")
    
    // Setup test user
    userID := createTestUser(t)
    
    // Create a test team to ensure there's at least one
    teamID := createTestTeam(t, userID)
    
    // Send request to get teams
    t.Log("Sending request to get teams")
    req, _ := http.NewRequest("GET", "/api/teams", nil)
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200", w.Code)
    }
    
    // Parse response to get teams
    var response struct {
        Teams []models.Team `json:"teams"`
    }
    
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    t.Logf("Retrieved %d teams successfully", len(response.Teams))
    
    if len(response.Teams) == 0 {
        t.Fatalf("Expected at least one team, got none")
    }
    
    // Verify the test team is in the list
    found := false
    for _, team := range response.Teams {
        if team.ID == teamID {
            found = true
            break
        }
    }
    
    if !found {
        t.Fatalf("Test team (ID: %s) not found in teams list", teamID)
    }
    
    t.Log("TestGetTeams passed")
}

func TestGetTeamDetails(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamDetails")
    
    // Setup test user and team
    userID := createTestUser(t)
    teamID := createTestTeam(t, userID)
    
    // Send request to get team details
    t.Log("Sending request to get team details")
    req, _ := http.NewRequest("GET", fmt.Sprintf("/api/team/%s", teamID), nil)
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200", w.Code)
    }
    
    // Parse response
    var team models.Team
    if err := json.Unmarshal(w.Body.Bytes(), &team); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    if team.ID != teamID {
        t.Fatalf("Retrieved team ID %s doesn't match expected ID %s", team.ID, teamID)
    }
    
    t.Logf("Successfully retrieved team details for team %s", teamID)
    t.Log("TestGetTeamDetails passed")
}

func TestJoinTeam(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestJoinTeam")
    
    // Create admin user and team
    adminID := createTestUser(t)
    teamID := createTestTeam(t, adminID)
    
    // Create a second user who will join the team
    secondUserID := createSecondTestUser(t)
    
    // Join team request
    joinRequest := struct {
        TeamID      string `json:"team_id"`
        TeamPassword string `json:"team_password"`
    }{
        TeamID:      teamID,
        TeamPassword: "testpassword", // Use the same password created with the team
    }
    
    jsonJoin, _ := json.Marshal(joinRequest)
    req, _ := http.NewRequest("POST", "/api/team/join", bytes.NewBuffer(jsonJoin))
    req.Header.Set("Content-Type", "application/json")
    
    // Add auth context for second user
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", secondUserID)
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200, response: %s", 
            w.Code, w.Body.String())
    }
    
    // Verify user joined by getting team members
    req, _ = http.NewRequest("GET", fmt.Sprintf("/api/team/%s/members", teamID), nil)
    ctx = req.Context()
    ctx = context.WithValue(ctx, "userID", adminID) // Admin can get members
    req = req.WithContext(ctx)
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    var members []models.TeamMemberDetails
    if err := json.Unmarshal(w.Body.Bytes(), &members); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    // Check if the second user is in the team
    secondUserFound := false
    for _, member := range members {
        if member.ID == secondUserID {
            secondUserFound = true
            break
        }
    }
    
    if !secondUserFound {
        t.Fatalf("Second user not found in team members after joining")
    }
    
    t.Log("TestJoinTeam passed")
}