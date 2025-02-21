package todos_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTodoQueries(DB *sql.DB) *db.Queries {
    return db.New(DB)
}

func NewTodoRepository(db *sql.DB) domain.TodoRepository {
    querier := NewTodoQueries(db)
    return &TodoRepository{querier: querier}
}