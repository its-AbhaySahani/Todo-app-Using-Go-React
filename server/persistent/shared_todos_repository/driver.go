package shared_todos_repository

import (
	"database/sql"

	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

func NewSharedTodoQueries(DB *sql.DB) *db.Queries {
	return db.New(DB)
}

func NewSharedTodoRepository(DB *sql.DB) domain.SharedTodoRepository {
	querier := NewSharedTodoQueries(DB)
	return &SharedTodoRepository{querier: querier}
}
