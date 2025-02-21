package team_members_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)



func NewTeamMemberRepository(DB *sql.DB) *TeamMemberRepository {
    querier := db.New(DB)
    return &TeamMemberRepository{querier: querier}
}