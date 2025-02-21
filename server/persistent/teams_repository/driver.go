package teams_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)


func NewTeamRepository(DB *sql.DB) *TeamRepository {
    querier := db.New(DB)
    return &TeamRepository{querier: querier}
}