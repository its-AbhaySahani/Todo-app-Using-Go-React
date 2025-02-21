package team_todos_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamTodoQueries(DB *sql.DB) *db.Queries {
    return db.New(DB)
}

func NewTeamTodoRepository(db *sql.DB) domain.TeamTodoRepository {
    querier := NewTeamTodoQueries(db)
    return &TeamTodoRepository{querier: querier}
}