package routine_repository

import (
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

func NewRoutineRepository(DB *sql.DB) *RoutineRepository {
    querier := db.New(DB)
    return &RoutineRepository{querier: querier}
}