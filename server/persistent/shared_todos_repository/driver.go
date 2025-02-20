package shared_todos_repository

import (
	"database/sql"

	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/server/models/db"
	
	//"github.com/its-AbhaySahani/Todo-app-Using-Go-React/server/persistent/models/db"
)

func NewSharedTodoQueries(db *sql.DB) *db.Queries {
	querier := db.New(db)
	return &sharedTodoRepository{querier: querier}

}
