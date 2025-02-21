package team_members_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamMemberQueries(DB *sql.DB) *db.Queries {
    return db.New(DB)
}

func NewTeamMemberRepository(db *sql.DB) domain.TeamMemberRepository {
    querier := NewTeamMemberQueries(db)
    return &TeamMemberRepository{querier: querier}
}