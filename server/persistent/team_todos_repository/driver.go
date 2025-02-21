package team_todos_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

func NewTeamTodoRepository(DB *sql.DB) *TeamTodoRepository {
    querier := db.New(DB)
    return &TeamTodoRepository{querier: querier}
}