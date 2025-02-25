package TestCases

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/middleware"
)

func TestAuthMiddlewareWithValidToken(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAuthMiddlewareWithValidToken")
    fmt.Println("Testing authentication middleware with valid token")
    
    // Generate a valid token
    token, err := generateToken("testuser", "test-user-id")
    if err != nil {
        t.Fatal("Failed to generate token:", err)
    }

    // Create a request with the token in the Authorization header
    req, err := http.NewRequest("GET", "/api/protected", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a test handler that will be called if auth passes
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check that userID was added to context
        userID, ok := r.Context().Value("userID").(string)
        if !ok {
            t.Error("userID not found in context")
        } else {
            fmt.Printf("userID in context: %s\n", userID)
        }
        
        // Write a success response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("protected"))
    })
    
    // Create the middleware handler
    handler := middleware.AuthMiddleware(testHandler)
    
    // Serve the request
    fmt.Println("Sending request with valid token")
    handler.ServeHTTP(rr, req)

    // Check the status code - should be OK
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
        fmt.Printf("Response body: %s\n", rr.Body.String())
    } else {
        fmt.Printf("Status code verified: %d (OK)\n", status)
        fmt.Printf("Response body: %s\n", rr.Body.String())
    }
    
    fmt.Println("Auth middleware with valid token test passed")
}

func TestAuthMiddlewareWithInvalidToken(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAuthMiddlewareWithInvalidToken")
    fmt.Println("Testing authentication middleware with invalid token")
    
    // Create a request with an invalid token
    req, err := http.NewRequest("GET", "/api/protected", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }
    req.Header.Set("Authorization", "Bearer invalidtoken")

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a test handler that should not be called
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Error("Handler was called with invalid token")
    })
    
    // Create the middleware handler
    handler := middleware.AuthMiddleware(testHandler)
    
    // Serve the request
    fmt.Println("Sending request with invalid token")
    handler.ServeHTTP(rr, req)

    // Check the status code - should be Bad Request or Unauthorized
    if status := rr.Code; status != http.StatusBadRequest && status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want either %v or %v", 
            status, http.StatusBadRequest, http.StatusUnauthorized)
    } else {
        // Fix ternary operator - Go doesn't have this feature
        var statusText string
        if status == http.StatusBadRequest {
            statusText = "Bad Request"
        } else {
            statusText = "Unauthorized"
        }
        fmt.Printf("Status code verified: %d (%s)\n", status, statusText)
    }
    
    fmt.Printf("Response body: %s\n", rr.Body.String())
    fmt.Println("Auth middleware with invalid token test passed")
}

func TestAuthMiddlewareWithNoToken(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestAuthMiddlewareWithNoToken")
    fmt.Println("Testing authentication middleware with no token")
    
    // Create a request with no token
    req, err := http.NewRequest("GET", "/api/protected", nil)
    if err != nil {
        t.Fatal("Failed to create request:", err)
    }

    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create a test handler that should not be called
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Error("Handler was called with no token")
    })
    
    // Create the middleware handler
    handler := middleware.AuthMiddleware(testHandler)
    
    // Serve the request
    fmt.Println("Sending request with no token")
    handler.ServeHTTP(rr, req)

    // Check the status code - should be unauthorized
    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
    } else {
        fmt.Printf("Status code verified: %d (Unauthorized)\n", status)
        fmt.Printf("Response body: %s\n", rr.Body.String())
    }
    
    fmt.Println("Auth middleware with no token test passed")
} 