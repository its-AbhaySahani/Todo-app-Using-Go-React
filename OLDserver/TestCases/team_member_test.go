// server/TestCases/team_member_test.go

package TestCases

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

var router *mux.Router

func init() {
    router = mux.NewRouter()
    router.HandleFunc("/api/register", middleware.Register).Methods("POST")
    router.HandleFunc("/api/login", middleware.Login).Methods("POST")
    
    apiRouter := router.PathPrefix("/api").Subrouter()
    apiRouter.Use(middleware.AuthMiddleware)
    
    // Team member routes
    apiRouter.HandleFunc("/team", middleware.CreateTeam).Methods("POST")
    apiRouter.HandleFunc("/teams", middleware.GetTeams).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/members", middleware.GetTeamMembers).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/member", middleware.AddTeamMember).Methods("POST")
    apiRouter.HandleFunc("/team/{teamId}/member/{userId}", middleware.RemoveTeamMember).Methods("DELETE")
}

// Helper functions
func createTestUser(t *testing.T) string {
    // Check if test user already exists
    testUserID := "todo-test-user-id"
    req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%s", testUserID), nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code == http.StatusOK {
        t.Logf("Test user already exists with ID: %s", testUserID)
        return testUserID
    }
    
    // Create test user if not exists
    user := models.User{
        Username: "testuser",
        Password: "testpassword",
    }
    
    jsonUser, _ := json.Marshal(user)
    req, _ = http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonUser))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Fatalf("Failed to create test user, status: %d", w.Code)
    }
    
    var response struct {
        ID string `json:"id"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    return response.ID
}

func createSecondTestUser(t *testing.T) string {
    user := models.User{
        Username: fmt.Sprintf("testuser2-%s", uuid.New().String()[:8]),
        Password: "testpassword",
    }
    
    jsonUser, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonUser))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %v want %v, response: %v", 
            w.Code, http.StatusOK, w.Body.String())
    }
    
    var response struct {
        ID string `json:"id"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    t.Logf("Created second test user with ID: %s", response.ID)
    return response.ID
}

func createTestTeam(t *testing.T, userID string) string {
    t.Log("Test team doesn't exist, creating new team")
    
    team := models.Team{
        Name:        "Test Team",
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
    
    if w.Code != http.StatusOK {
        t.Fatalf("Failed to create team: status %d, response: %s", w.Code, w.Body.String())
    }
    
    var response struct {
        ID string `json:"id"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse team creation response: %v", err)
    }
    
    if response.ID == "" {
        t.Fatalf("Team created but no ID returned")
    }
    
    t.Logf("Created test team with ID: %s", response.ID)
    return response.ID
}

func TestGetTeamMembers(t *testing.T) {
    fmt.Println("=== RUN   TestGetTeamMembers")
    
    // Create test user
    userID := createTestUser(t)
    
    // Create test team
    teamID := createTestTeam(t, userID)
    
    // Send request to get team members
    req, _ := http.NewRequest("GET", fmt.Sprintf("/api/team/%s/members", teamID), nil)
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %v want 200", w.Code)
    }
    
    // Parse response
    var members []models.TeamMemberDetails
    if err := json.Unmarshal(w.Body.Bytes(), &members); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    // The team creator (admin) should be a member by default
    if len(members) < 1 {
        t.Fatalf("Expected at least one team member, got none")
    }
    
    // Check if the admin is in the members list
    adminFound := false
    for _, member := range members {
        if member.ID == userID && member.IsAdmin {
            adminFound = true
            break
        }
    }
    
    if !adminFound {
        t.Fatalf("Admin not found in team members list")
    }
    
    t.Log("TestGetTeamMembers passed")
}

func TestAddTeamMember(t *testing.T) {
    fmt.Println("=== RUN   TestAddTeamMember")
    
    // Create test user (admin)
    userID := createTestUser(t)
    
    // Create test team
    teamID := createTestTeam(t, userID)
    
    // Create second test user to add as a member
    secondUserID := createSecondTestUser(t)
    
    // Prepare request body
    teamMember := struct {
        UserID  string `json:"user_id"`
        IsAdmin bool   `json:"is_admin"`
    }{
        UserID:  secondUserID,
        IsAdmin: false,
    }
    
    jsonMember, _ := json.Marshal(teamMember)
    req, _ := http.NewRequest("POST", fmt.Sprintf("/api/team/%s/member", teamID), bytes.NewBuffer(jsonMember))
    req.Header.Set("Content-Type", "application/json")
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID) // Admin user
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200. Body: %s", 
            w.Code, w.Body.String())
    }
    
    // Verify the member was added by getting team members
    req, _ = http.NewRequest("GET", fmt.Sprintf("/api/team/%s/members", teamID), nil)
    ctx = req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Fatalf("Failed to get team members: %d", w.Code)
    }
    
    var members []models.TeamMemberDetails
    if err := json.Unmarshal(w.Body.Bytes(), &members); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    // Check if the new member is in the list
    memberFound := false
    for _, member := range members {
        if member.ID == secondUserID {
            memberFound = true
            break
        }
    }
    
    if !memberFound {
        t.Fatalf("Added member not found in team members list")
    }
    
    t.Log("TestAddTeamMember passed")
}

func TestRemoveTeamMember(t *testing.T) {
    fmt.Println("=== RUN   TestRemoveTeamMember")
    
    // Create test user (admin)
    userID := createTestUser(t)
    
    // Create test team
    teamID := createTestTeam(t, userID)
    
    // Create second test user to add and then remove
    secondUserID := createSecondTestUser(t)
    
    // Add the second user as a member first
    teamMember := struct {
        UserID  string `json:"user_id"`
        IsAdmin bool   `json:"is_admin"`
    }{
        UserID:  secondUserID,
        IsAdmin: false,
    }
    
    jsonMember, _ := json.Marshal(teamMember)
    req, _ := http.NewRequest("POST", fmt.Sprintf("/api/team/%s/member", teamID), bytes.NewBuffer(jsonMember))
    req.Header.Set("Content-Type", "application/json")
    
    // Add auth context
    ctx := req.Context()
    ctx = context.WithValue(ctx, "userID", userID) // Admin user
    req = req.WithContext(ctx)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Fatalf("Failed to add member: %d, %s", w.Code, w.Body.String())
    }
    
    // Now remove the member
    req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/team/%s/member/%s", teamID, secondUserID), nil)
    ctx = req.Context()
    ctx = context.WithValue(ctx, "userID", userID) // Admin user
    req = req.WithContext(ctx)
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Check response
    if w.Code != http.StatusOK {
        t.Fatalf("Handler returned wrong status code: got %d want 200, response: %s", 
            w.Code, w.Body.String())
    }
    
    // Verify the member was removed by getting team members
    req, _ = http.NewRequest("GET", fmt.Sprintf("/api/team/%s/members", teamID), nil)
    ctx = req.Context()
    ctx = context.WithValue(ctx, "userID", userID)
    req = req.WithContext(ctx)
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    var members []models.TeamMemberDetails
    if err := json.Unmarshal(w.Body.Bytes(), &members); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
    
    // Check that the removed member is not in the list
    for _, member := range members {
        if member.ID == secondUserID {
            t.Fatalf("Removed member still found in team members list")
        }
    }
    
    t.Log("TestRemoveTeamMember passed")
}