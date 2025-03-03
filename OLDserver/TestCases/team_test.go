package TestCases

import (
    "fmt"
    "testing"


    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "golang.org/x/crypto/bcrypt"
)

// Constants for team tests
const testTeamAdminID = "team-test-admin-id"
const testTeamAdminUsername = "team_test_admin"
var testTeamID string // Will be populated when we create a test team

// Helper function to ensure the admin test user exists
func ensureTeamAdminExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testTeamAdminID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if team admin exists: %v", err)
    }

    if count == 0 {
        // Create the admin user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testTeamAdminID, testTeamAdminUsername, "$2a$10$TestHashedPasswordForTeamAdmin",
        )
        if err != nil {
            t.Fatalf("Failed to create team admin: %v", err)
        }
        fmt.Println("Created team admin with ID:", testTeamAdminID)
    } else {
        fmt.Println("Team admin already exists with ID:", testTeamAdminID)
    }
}

// Helper function to cleanup team test data
func cleanupTeamTestData() {
    // Delete team members
    if testTeamID != "" {
        _, err := database.DB.Exec("DELETE FROM team_members WHERE team_id = ?", testTeamID)
        if err != nil {
            fmt.Println("Error cleaning up team members:", err)
        } else {
            fmt.Println("Team members cleaned up")
        }
        
        // Delete the team
        _, err = database.DB.Exec("DELETE FROM teams WHERE id = ?", testTeamID)
        if err != nil {
            fmt.Println("Error cleaning up team:", err)
        } else {
            fmt.Println("Team cleaned up")
        }
    }
    
    // Delete test user
    _, err := database.DB.Exec("DELETE FROM users WHERE id = ?", testTeamAdminID)
    if err != nil {
        fmt.Println("Error cleaning up team admin:", err)
    } else {
        fmt.Println("Team admin cleaned up")
    }
    
    testTeamID = "" // Reset the team ID
}

// TestCreateTeam tests the CreateTeam function directly
func TestCreateTeam(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTeam")
    fmt.Println("Testing creating a new team")
    
    // Clean up before test to ensure fresh state
    cleanupTeamTestData()
    
    // Ensure admin user exists
    ensureTeamAdminExists(t)
    
    // Create a new team directly using the model function
    name := "Functional Test Team"
    password := "testteampassword"
    
    // Hash the password as it would be in real application
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    team, err := models.CreateTeam(name, string(hashedPassword), testTeamAdminID)
    if err != nil {
        t.Fatalf("Failed to create team: %v", err)
    }
    
    // Verify the team was created with the correct data
    if team.ID == "" {
        t.Errorf("Team was created but ID is empty")
    } else {
        fmt.Printf("Created team with ID: %s\n", team.ID)
    }
    
    if team.Name != name {
        t.Errorf("Expected team name '%s', got '%s'", name, team.Name)
    }
    
    if team.AdminID != testTeamAdminID {
        t.Errorf("Expected admin ID '%s', got '%s'", testTeamAdminID, team.AdminID)
    }
    
    // Verify the team exists in the database
    var dbName, dbAdminID string
    
    err = database.DB.QueryRow(
        "SELECT name, admin_id FROM teams WHERE id = ?", 
        team.ID,
    ).Scan(&dbName, &dbAdminID)
    
    if err != nil {
        t.Fatalf("Failed to retrieve team from database: %v", err)
    }
    
    if dbName != name {
        t.Errorf("Database team name '%s' doesn't match expected name '%s'", dbName, name)
    }
    
    if dbAdminID != testTeamAdminID {
        t.Errorf("Database admin ID '%s' doesn't match expected admin ID '%s'", dbAdminID, testTeamAdminID)
    }
    
    // Save ID for later tests
    testTeamID = team.ID
    
    fmt.Println("CreateTeam test passed")
}

// TestGetTeams tests the GetTeams function directly
func TestGetTeams(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeams")
    fmt.Println("Testing getting teams for a user")
    
    // Ensure we have a team to retrieve
    if testTeamID == "" {
        TestCreateTeam(t)
        if testTeamID == "" {
            t.Fatal("Failed to create a team for retrieval test")
        }
    }
    
    // Get teams for the admin
    teams, err := models.GetTeams(testTeamAdminID)
    if err != nil {
        t.Fatalf("Failed to get teams: %v", err)
    }
    
    // Verify we got at least one team
    if len(teams) == 0 {
        t.Errorf("Expected at least one team, got none")
    } else {
        fmt.Printf("Retrieved %d teams\n", len(teams))
        
        // Find our test team
        found := false
        for _, team := range teams {
            if team.ID == testTeamID {
                found = true
                if team.Name != "Functional Test Team" {
                    t.Errorf("Expected team name 'Functional Test Team', got '%s'", team.Name)
                }
                
                if team.AdminID != testTeamAdminID {
                    t.Errorf("Expected admin ID '%s', got '%s'", testTeamAdminID, team.AdminID)
                }
            }
        }
        
        if !found {
            t.Errorf("Test team not found in retrieved teams")
        }
    }
    
    fmt.Println("GetTeams test passed")
}

// TestGetTeamByID tests the GetTeamByID function directly
func TestGetTeamByID(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamByID")
    fmt.Println("Testing getting a team by ID")
    
    // Ensure we have a team to retrieve
    if testTeamID == "" {
        TestCreateTeam(t)
        if testTeamID == "" {
            t.Fatal("Failed to create a team for retrieval test")
        }
    }
    
    // Get the team by ID
    team, err := models.GetTeamByID(testTeamID)
    if err != nil {
        t.Fatalf("Failed to get team by ID: %v", err)
    }
    
    // Verify we got the correct team
    if team.ID != testTeamID {
        t.Errorf("Expected team ID '%s', got '%s'", testTeamID, team.ID)
    }
    
    if team.Name != "Functional Test Team" {
        t.Errorf("Expected team name 'Functional Test Team', got '%s'", team.Name)
    }
    
    if team.AdminID != testTeamAdminID {
        t.Errorf("Expected admin ID '%s', got '%s'", testTeamAdminID, team.AdminID)
    }
    
    fmt.Println("GetTeamByID test passed")
    
    // Clean up after team tests
    cleanupTeamTestData()
}

