package teams_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamQueries(DB *sql.DB) *db.Queries {
    return db.New(DB)
}

func NewTeamRepository(db *sql.DB) domain.TeamRepository {
    querier := NewTeamQueries(db)
    return &TeamRepository{querier: querier}
}