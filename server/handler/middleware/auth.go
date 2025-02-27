package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

type contextKey string

const (
    UserIDKey contextKey = "userID"
)

var jwtKey = []byte("ZLR+ZInOHXQst1seVlV6JVuZe1k3vasV1BRyqAHAyaY=") // This should match your key in api.go

type Claims struct {
    Username string `json:"username"`
    UserID   string `json:"user_id"`
    jwt.StandardClaims
}

// AuthMiddleware verifies JWT tokens and adds user ID to context
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        // Format should be "Bearer token"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
            return
        }

        tokenStr := parts[1]

        // Parse and validate token
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        // Add userID to request context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}