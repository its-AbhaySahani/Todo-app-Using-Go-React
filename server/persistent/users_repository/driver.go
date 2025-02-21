package users_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)


func NewUserRepository(DB *sql.DB) *UserRepository {
    querier := db.New(DB)
    return &UserRepository{querier: querier}
}