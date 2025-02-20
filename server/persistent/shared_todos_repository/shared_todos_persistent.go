package shared_todos_repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type SharedTodoRepository struct {
	querier *db.Queries
}

func (r *SharedTodoRepository) CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) error {
	id := uuid.New().String()
	date := time.Now()
	timeNow := time.Now()
	return r.querier.CreateSharedTodo(ctx, db.CreateSharedTodoParams{
		ID:          id,
		Task:        sql.NullString{String: task, Valid: true},
		Description: sql.NullString{String: description, Valid: true},
		Done:        sql.NullBool{Bool: done, Valid: true},
		Important:   sql.NullBool{Bool: important, Valid: true},
		UserID:      sql.NullString{String: userID, Valid: true},
		Date:        sql.NullTime{Time: date, Valid: true},
		Time:        sql.NullTime{Time: timeNow, Valid: true},
		SharedBy:    sql.NullString{String: sharedBy, Valid: true},
	})
}

func (r *SharedTodoRepository) GetSharedTodos(ctx context.Context, userID string) ([]domain.SharedTodo, error) {
	todos, err := r.querier.GetSharedTodos(ctx, sql.NullString{String: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	var result []domain.SharedTodo
	for _, todo := range todos {
		result = append(result, domain.SharedTodo{
			ID:          todo.ID,
			Task:        todo.Task.String,
			Description: todo.Description.String,
			Done:        todo.Done.Bool,
			Important:   todo.Important.Bool,
			UserID:      todo.UserID.String,
			Date:        todo.Date.Time,
			Time:        todo.Time.Time,
			SharedBy:    todo.SharedBy.String,
		})
	}
	return result, nil
}

func (r *SharedTodoRepository) GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]domain.SharedTodo, error) {
	todos, err := r.querier.GetSharedByMeTodos(ctx, sql.NullString{String: sharedBy, Valid: true})
	if err != nil {
		return nil, err
	}
	var result []domain.SharedTodo
	for _, todo := range todos {
		result = append(result, domain.SharedTodo{
			ID:          todo.ID,
			Task:        todo.Task.String,
			Description: todo.Description.String,
			Done:        todo.Done.Bool,
			Important:   todo.Important.Bool,
			UserID:      todo.UserID.String,
			Date:        todo.Date.Time,
			Time:        todo.Time.Time,
			SharedBy:    todo.SharedBy.String,
		})
	}
	return result, nil
}
