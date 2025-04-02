package helpers

// Auth request/response types
type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterResponse struct {
    ID string `json:"id"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    ID    string `json:"id"`
}