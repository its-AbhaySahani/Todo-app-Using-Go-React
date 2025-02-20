package models

import (
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func CreateUser(username, password string) (User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return User{}, err
    }

    id := uuid.New().String()
    _, err = database.DB.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)", id, username, string(hashedPassword))
    if err != nil {
        return User{}, err
    }

    return User{ID: id, Username: username, Password: string(hashedPassword)}, nil
}

func GetUserByUsername(username string) (User, error) {
    var user User
    err := database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return User{}, err
    }
    return user, nil
}

func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}