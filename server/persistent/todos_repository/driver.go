package todos_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

func NewTodoRepository(DB *sql.DB) *TodoRepository {
    querier := db.New(DB)
    return &TodoRepository{
        querier: querier,
        db:      DB,  // Store DB connection
    }
}