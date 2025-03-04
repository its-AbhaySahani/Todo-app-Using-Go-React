package shared_todos_repository

import (
	"database/sql"

	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

func NewSharedTodoRepository(DB *sql.DB) *SharedTodoRepository {
    querier := db.New(DB)
    return &SharedTodoRepository{
        querier: querier,
        db:      DB,  // Store DB connection
    }
}