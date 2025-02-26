package TestCases

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
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

// TestRegister tests the user registration endpoint
func TestRegister(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestRegister")
    fmt.Println("Testing user registration with username 'testuser'")
    
    // Clean up before test
    cleanupTestData()
    
    // Create a request body with test credentials
    var jsonStr = []byte(`{"username":"testuser", "password":"testpassword"}`)
    req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create handler
    handler := http.HandlerFunc(middleware.Register)
    
    // Serve the request
    fmt.Println("Sending registration request")
    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    }

    // Parse the response
    var response models.User
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate the response
    fmt.Printf("Registration response: %+v\n", response)
    
    if response.Username != "testuser" {
        t.Errorf("Handler returned unexpected username: got %v want %v", response.Username, "testuser")
    } else {
        fmt.Println("Username verified: testuser")
    }
    
    if response.ID == "" {
        t.Errorf("Handler returned empty user ID")
    } else {
        fmt.Printf("User ID verified: %s\n", response.ID)
    }
    
    // Password should be hashed, not returned as plaintext
    if response.Password == "testpassword" {
        t.Errorf("Password was returned as plaintext")
    } else {
        fmt.Println("Password is properly hashed")
    }
    
    fmt.Println("Registration test passed")
}

// TestLoginWithValidCredentials tests the login endpoint with valid credentials
func TestLoginWithValidCredentials(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestLoginWithValidCredentials")
    fmt.Println("Testing login with valid credentials")
    

	// Clean up before test to ensure fresh state
    cleanupTestData()

    // First, ensure the test user exists
    user, err := models.CreateUser("testuser", "testpassword")
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    fmt.Printf("Test user created with ID: %s\n", user.ID)

    // Create a request body with valid credentials
    var jsonStr = []byte(`{"username":"testuser", "password":"testpassword"}`)
    req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create handler
    handler := http.HandlerFunc(middleware.Login)
    
    // Serve the request
    fmt.Println("Sending login request with valid credentials")
    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
        return
    } else {
        fmt.Printf("Status code verified: %d (OK)\n", status)
    }

    // Parse the response
    var response map[string]interface{}
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v\nResponse body: %s", err, rr.Body.String())
    }

    // Validate that a token was returned
    fmt.Printf("Login response: %v\n", response)
    
    token, exists := response["token"]
    if !exists {
        t.Errorf("Handler did not return a token field")
    } else if token == "" {
        t.Errorf("Handler returned an empty token")
    } else {
        fmt.Printf("Token verified (length: %d characters)\n", len(token.(string)))
    }

    // Check for token cookie
    cookies := rr.Result().Cookies()
    foundTokenCookie := false
    for _, cookie := range cookies {
        if cookie.Name == "token" {
            foundTokenCookie = true
            fmt.Printf("Token cookie found (expires: %v)\n", cookie.Expires)
            break
        }
    }
    if !foundTokenCookie {
        t.Errorf("Handler did not set token cookie")
    }
    
    fmt.Println("Login with valid credentials test passed")
}

// TestLoginWithInvalidCredentials tests the login endpoint with invalid credentials
func TestLoginWithInvalidCredentials(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestLoginWithInvalidCredentials")
    fmt.Println("Testing login with invalid credentials (wrong password)")
    
	cleanupTestData()
    // Ensure test user exists with correct password
    user, err := models.CreateUser("testuser", "testpassword")
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    fmt.Printf("Test user created with ID: %s\n", user.ID)

    // Create a request body with invalid password
    var jsonStr = []byte(`{"username":"testuser", "password":"wrongpassword"}`)
    req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create handler
    handler := http.HandlerFunc(middleware.Login)
    
    // Serve the request
    fmt.Println("Sending login request with wrong password")
    handler.ServeHTTP(rr, req)

    // Check the status code - should be unauthorized
    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    } else {
        fmt.Printf("Status code verified: %d (Unauthorized)\n", status)
    }
    
    // Show the error message
    fmt.Printf("Error message: %s\n", rr.Body.String())
    
    fmt.Println("Login with invalid credentials test passed")
}

// TestLoginWithNonExistentUser tests login with a username that doesn't exist
func TestLoginWithNonExistentUser(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestLoginWithNonExistentUser")
    fmt.Println("Testing login with non-existent user")

    // Create a request body with non-existent user
    var jsonStr = []byte(`{"username":"nonexistentuser", "password":"anypassword"}`)
    req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create handler
    handler := http.HandlerFunc(middleware.Login)
    
    // Serve the request
    fmt.Println("Sending login request with non-existent user")
    handler.ServeHTTP(rr, req)

    // Check the status code - should be unauthorized
    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    } else {
        fmt.Printf("Status code verified: %d (Unauthorized)\n", status)
    }
    
    // Show the error message
    fmt.Printf("Error message: %s\n", rr.Body.String())
    
    fmt.Println("Login with non-existent user test passed")
}