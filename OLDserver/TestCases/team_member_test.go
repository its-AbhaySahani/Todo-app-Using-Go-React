package TestCases

import (
    "fmt"
    "testing"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
)

// Constants for team member tests
const testTeamOwnerID = "team-member-owner-id"
const testTeamOwnerUsername = "team_member_owner"
const testTeamMemberID = "team-member-user-id"
const testTeamMemberUsername = "team_member_user"
var testTeamMemberTeamID string // Will be populated when we create a test team

// Helper function to ensure the team owner exists
func ensureTeamOwnerExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testTeamOwnerID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if team owner exists: %v", err)
    }

    if count == 0 {
        // Create the owner user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testTeamOwnerID, testTeamOwnerUsername, "$2a$10$TestHashedPasswordForTeamOwner",
        )
        if err != nil {
            t.Fatalf("Failed to create team owner: %v", err)
        }
        fmt.Println("Created team owner with ID:", testTeamOwnerID)
    } else {
        fmt.Println("Team owner already exists with ID:", testTeamOwnerID)
    }
}

// Helper function to ensure the team member user exists
func ensureTeamMemberExists(t *testing.T) {
    // Check if the user already exists
    var count int
    err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", testTeamMemberID).Scan(&count)
    if err != nil {
        t.Fatalf("Failed to check if team member exists: %v", err)
    }

    if count == 0 {
        // Create the member user if it doesn't exist
        _, err := database.DB.Exec(
            "INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
            testTeamMemberID, testTeamMemberUsername, "$2a$10$TestHashedPasswordForTeamMember",
        )
        if err != nil {
            t.Fatalf("Failed to create team member: %v", err)
        }
        fmt.Println("Created team member with ID:", testTeamMemberID)
    } else {
        fmt.Println("Team member already exists with ID:", testTeamMemberID)
    }
}

// Helper function to cleanup team member test data
func cleanupTeamMemberTestData() {
    // Delete team members
    if testTeamMemberTeamID != "" {
        _, err := database.DB.Exec("DELETE FROM team_members WHERE team_id = ?", testTeamMemberTeamID)
        if err != nil {
            fmt.Println("Error cleaning up team members:", err)
        } else {
            fmt.Println("Team members cleaned up")
        }
        
        // Delete the team
        _, err = database.DB.Exec("DELETE FROM teams WHERE id = ?", testTeamMemberTeamID)
        if err != nil {
            fmt.Println("Error cleaning up team:", err)
        } else {
            fmt.Println("Team cleaned up")
        }
    }
    
    // Delete test users
    _, err := database.DB.Exec("DELETE FROM users WHERE id IN (?, ?)", testTeamOwnerID, testTeamMemberID)
    if err != nil {
        fmt.Println("Error cleaning up test users:", err)
    } else {
        fmt.Println("Test users cleaned up")
    }
    
    testTeamMemberTeamID = "" // Reset the team ID
}

// Helper function to create a test team
func createTestTeamForMemberTests(t *testing.T) string {
    // Check if we already have a team ID
    if testTeamMemberTeamID != "" {
        var count int
        err := database.DB.QueryRow("SELECT COUNT(*) FROM teams WHERE id = ?", testTeamMemberTeamID).Scan(&count)
        if err == nil && count > 0 {
            fmt.Println("Using existing test team with ID:", testTeamMemberTeamID)
            return testTeamMemberTeamID
        }
    }
    
    // Ensure owner user exists
    ensureTeamOwnerExists(t)
    
    // Create a new team directly using the model function
    name := fmt.Sprintf("Team Member Test Team %s", uuid.New().String()[:8])
    password := "testteampassword"
    
    team, err := models.CreateTeam(name, password, testTeamOwnerID)
    if err != nil {
        t.Fatalf("Failed to create team for member tests: %v", err)
    }
    
    fmt.Printf("Created team with ID: %s for member tests\n", team.ID)
    testTeamMemberTeamID = team.ID
    return team.ID
}

// TestAddTeamMember tests the AddTeamMember function directly
func TestAddTeamMember(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAddTeamMember")
    fmt.Println("Testing adding a member to a team")
    
    // Clean up before test to ensure fresh state
    cleanupTeamMemberTestData()
    
    // Create a test team and ensure member user exists
    teamID := createTestTeamForMemberTests(t)
    ensureTeamMemberExists(t)
    
    // Add the member to the team directly using the model function
    err := models.AddTeamMember(teamID, testTeamMemberUsername, testTeamOwnerID)
    if err != nil {
        t.Fatalf("Failed to add team member: %v", err)
    }
    
    // Verify the member was added by checking the database
    var count int
    err = database.DB.QueryRow(
        "SELECT COUNT(*) FROM team_members WHERE team_id = ? AND user_id = ?", 
        teamID, testTeamMemberID,
    ).Scan(&count)
    
    if err != nil {
        t.Fatalf("Failed to check if team member was added: %v", err)
    }
    
    if count != 1 {
        t.Errorf("Expected team member to be added, but it wasn't found in the database")
    } else {
        fmt.Println("Team member was successfully added to the database")
    }
    
    fmt.Println("AddTeamMember test passed")
}

// TestGetTeamMembers tests the GetTeamMembers function directly
func TestGetTeamMembers(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamMembers")
    fmt.Println("Testing getting all members of a team")
    
    // First add a team member to ensure there's data to retrieve
    if testTeamMemberTeamID == "" {
        TestAddTeamMember(t)
        if testTeamMemberTeamID == "" {
            t.Fatal("Failed to create a team for team member retrieval test")
        }
    }
    
    // Get team members
    members, err := models.GetTeamMembers(testTeamMemberTeamID)
    if err != nil {
        t.Fatalf("Failed to get team members: %v", err)
    }
    
    // Verify we got at least two members (the owner and our added member)
    if len(members) < 2 {
        t.Errorf("Expected at least 2 team members, got %d", len(members))
    } else {
        fmt.Printf("Retrieved %d team members\n", len(members))
        
        // Find our added member
        foundMember := false
        foundOwner := false
        
        for _, member := range members {
            if member.ID == testTeamMemberID {
                foundMember = true
                if member.IsAdmin {
                    t.Errorf("Expected team member to not be admin, but it is")
                }
            } else if member.ID == testTeamOwnerID {
                foundOwner = true
                if !member.IsAdmin {
                    t.Errorf("Expected team owner to be admin, but it's not")
                }
            }
        }
        
        if !foundMember {
            t.Errorf("Added team member not found in retrieved members")
        }
        
        if !foundOwner {
            t.Errorf("Team owner not found in retrieved members")
        }
    }
    
    fmt.Println("GetTeamMembers test passed")
}

// TestRemoveTeamMember tests the RemoveTeamMember function directly
func TestRemoveTeamMember(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestRemoveTeamMember")
    fmt.Println("Testing removing a member from a team")
    
    // First add a team member to ensure there's data to remove
    if testTeamMemberTeamID == "" {
        TestAddTeamMember(t)
        if testTeamMemberTeamID == "" {
            t.Fatal("Failed to create a team for team member removal test")
        }
    }
    
    // Remove the team member directly using the model function
    err := models.RemoveTeamMember(testTeamMemberTeamID, testTeamMemberID)
    if err != nil {
        t.Fatalf("Failed to remove team member: %v", err)
    }
    
    // Verify the member was removed by checking the database
    var count int
    err = database.DB.QueryRow(
        "SELECT COUNT(*) FROM team_members WHERE team_id = ? AND user_id = ?", 
        testTeamMemberTeamID, testTeamMemberID,
    ).Scan(&count)
    
    if err != nil {
        t.Fatalf("Failed to check if team member was removed: %v", err)
    }
    
    if count != 0 {
        t.Errorf("Expected team member to be removed, but it's still in the database")
    } else {
        fmt.Println("Team member was successfully removed from the database")
    }
    
    fmt.Println("RemoveTeamMember test passed")
    
    // Clean up after all team member tests
    cleanupTeamMemberTestData()
}

