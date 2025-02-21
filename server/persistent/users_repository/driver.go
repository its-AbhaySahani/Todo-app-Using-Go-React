package users_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewUserQueries(DB *sql.DB) *db.Queries {
    return db.New(DB)
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
    querier := NewUserQueries(db)
    return &UserRepository{querier: querier}
}