package TestCases

import (
    "fmt"
    "testing"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("ZLR+ZInOHXQst1seVlV6JVuZe1k3vasV1BRyqAHAyaY=")

func cleanupTestData() {
    // Delete test user if exists
    _, err := database.DB.Exec("DELETE FROM users WHERE username = ?", "testuser")
    if err != nil {
        fmt.Println("Error cleaning up test data:", err)
    } else {
        fmt.Println("Test data cleaned up")
    }
    
    // Delete any additional test users
    _, err = database.DB.Exec("DELETE FROM users WHERE username = ?", "testuser2")
    if err != nil {
        fmt.Println("Error cleaning up additional test data:", err)
    }
}

// Helper function to generate a JWT token for testing
func generateToken(username, userID string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &middleware.Claims{
        Username: username,
        UserID:   userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// TestCreateUser tests the user creation function directly
func TestCreateUser(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateUser")
    fmt.Println("Testing user creation with username 'testuser'")
    
    // Clean up before test
    cleanupTestData()
    
    // Call the CreateUser function directly
    username := "testuser"
    password := "testpassword"
    
    user, err := models.CreateUser(username, password)
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }
    
    // Validate the user object
    fmt.Printf("Created user: %+v\n", user)
    
    if user.Username != username {
        t.Errorf("Expected username %s, got %s", username, user.Username)
    }
    
    if user.ID == "" {
        t.Errorf("Expected a non-empty user ID")
    } else {
        fmt.Printf("User ID verified: %s\n", user.ID)
    }
    
    // Verify the password was hashed
    if user.Password == password {
        t.Errorf("Password was not hashed")
    } else {
        // Verify we can validate the password using bcrypt
        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
        if err != nil {
            t.Errorf("Password hash verification failed: %v", err)
        } else {
            fmt.Println("Password hash verification successful")
        }
    }
    
    // Verify the user exists in the database
    var dbUser models.User
    err = database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
        Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
    
    if err != nil {
        t.Fatalf("Failed to retrieve user from database: %v", err)
    }
    
    if dbUser.ID != user.ID {
        t.Errorf("Database user ID %s doesn't match created user ID %s", dbUser.ID, user.ID)
    }
    
    fmt.Println("User creation test passed")
}

// TestGetUserByUsername tests the GetUserByUsername function
func TestGetUserByUsername(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetUserByUsername")
    fmt.Println("Testing retrieving a user by username")
    
    // Clean up before test
    cleanupTestData()
    
    // Create a test user first
    username := "testuser"
    password := "testpassword"
    
    createdUser, err := models.CreateUser(username, password)
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    fmt.Printf("Created test user with ID: %s\n", createdUser.ID)
    
    // Now retrieve the user using the function
    user, err := models.GetUserByUsername(username)
    if err != nil {
        t.Fatalf("Failed to retrieve user: %v", err)
    }
    
    // Validate retrieved user
    if user.Username != username {
        t.Errorf("Expected username %s, got %s", username, user.Username)
    }
    
    if user.ID != createdUser.ID {
        t.Errorf("Expected user ID %s, got %s", createdUser.ID, user.ID)
    }
    
    fmt.Println("GetUserByUsername test passed")
    
    // Test non-existent user
    fmt.Println("Testing retrieving a non-existent user")
    _, err = models.GetUserByUsername("nonexistentuser")
    if err == nil {
        t.Errorf("Expected error when retrieving non-existent user, got nil")
    } else {
        fmt.Printf("Correctly got error for non-existent user: %v\n", err)
    }
}

// TestVerifyPassword tests the password verification function
func TestVerifyPassword(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestVerifyPassword")
    fmt.Println("Testing password verification")
    
    // Create a hashed password
    password := "testpassword"
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    // Test correct password
    err = models.VerifyPassword(string(hashedPassword), password)
    if err != nil {
        t.Errorf("Failed to verify correct password: %v", err)
    } else {
        fmt.Println("Password verification successful for correct password")
    }
    
    // Test incorrect password
    err = models.VerifyPassword(string(hashedPassword), "wrongpassword")
    if err == nil {
        t.Errorf("Expected error for incorrect password, got nil")
    } else {
        fmt.Printf("Correctly got error for incorrect password: %v\n", err)
    }
}

// TestDuplicateUsername tests creating a user with a duplicate username
func TestDuplicateUsername(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDuplicateUsername")
    fmt.Println("Testing creation of user with duplicate username")
    
    // Clean up before test
    cleanupTestData()
    
    // Create the first user
    username := "testuser"
    password := "testpassword"
    
    _, err := models.CreateUser(username, password)
    if err != nil {
        t.Fatalf("Failed to create first test user: %v", err)
    }
    
    // Try to create a second user with the same username
    _, err = models.CreateUser(username, "differentpassword")
    if err == nil {
        t.Errorf("Expected error when creating duplicate username, got nil")
    } else {
        fmt.Printf("Correctly got error for duplicate username: %v\n", err)
    }
    
    fmt.Println("Duplicate username test passed")
}

// TestTokenGeneration tests JWT token generation and validation
func TestTokenGeneration(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestTokenGeneration")
    fmt.Println("Testing JWT token generation and validation")
    
    username := "testuser"
    userID := "test-user-id"
    
    // Generate a token
    tokenString, err := generateToken(username, userID)
    if err != nil {
        t.Fatalf("Failed to generate token: %v", err)
    }
    
    if tokenString == "" {
        t.Fatalf("Generated token is empty")
    }
    
    fmt.Printf("Generated token (length: %d)\n", len(tokenString))
    
    // Parse and validate the token
    token, err := jwt.ParseWithClaims(tokenString, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    
    if err != nil {
        t.Fatalf("Failed to parse token: %v", err)
    }
    
    if !token.Valid {
        t.Fatalf("Token is not valid")
    }
    
    claims, ok := token.Claims.(*middleware.Claims)
    if !ok {
        t.Fatalf("Couldn't parse claims")
    }
    
    if claims.Username != username {
        t.Errorf("Expected username %s in token, got %s", username, claims.Username)
    }
    
    if claims.UserID != userID {
        t.Errorf("Expected userID %s in token, got %s", userID, claims.UserID)
    }
    
    fmt.Println("Token generation and validation test passed")
}

// TestCreateMultipleUsers tests creating multiple users
func TestCreateMultipleUsers(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateMultipleUsers")
    fmt.Println("Testing creation of multiple users")
    
    // Clean up before test
    cleanupTestData()
    
    // Create first user
    username1 := "testuser"
    password1 := "testpassword"
    
    user1, err := models.CreateUser(username1, password1)
    if err != nil {
        t.Fatalf("Failed to create first user: %v", err)
    }
    
    // Create second user
    username2 := "testuser2"
    password2 := "testpassword2"
    
    user2, err := models.CreateUser(username2, password2)
    if err != nil {
        t.Fatalf("Failed to create second user: %v", err)
    }
    
    // Verify both users have different IDs
    if user1.ID == user2.ID {
        t.Errorf("Both users have the same ID: %s", user1.ID)
    } else {
        fmt.Printf("User IDs are different: %s and %s\n", user1.ID, user2.ID)
    }
    
    // Retrieve both users to verify they exist in the database
    retrievedUser1, err := models.GetUserByUsername(username1)
    if err != nil {
        t.Fatalf("Failed to retrieve first user: %v", err)
    }
    
    retrievedUser2, err := models.GetUserByUsername(username2)
    if err != nil {
        t.Fatalf("Failed to retrieve second user: %v", err)
    }
    
    if retrievedUser1.ID != user1.ID {
        t.Errorf("Retrieved user1 ID %s doesn't match created user1 ID %s", retrievedUser1.ID, user1.ID)
    }
    
    if retrievedUser2.ID != user2.ID {
        t.Errorf("Retrieved user2 ID %s doesn't match created user2 ID %s", retrievedUser2.ID, user2.ID)
    }
    
    fmt.Println("Multiple users creation test passed")
}
