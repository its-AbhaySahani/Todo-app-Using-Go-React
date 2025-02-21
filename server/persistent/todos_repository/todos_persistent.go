package todos_repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type TodoRepository struct {
    querier *db.Queries
}

func (r *TodoRepository) CreateTodo(ctx context.Context, task, description string, done, important bool, userID string) error {
    id := uuid.New().String()
    date := time.Now()
    timeNow := time.Now()
    return r.querier.CreateTodo(ctx, db.CreateTodoParams{
        ID:          id,
        Task:        task,
        Description: sql.NullString{String: description, Valid: true},
        Done:        done,
        Important:   important,
        UserID:      sql.NullString{String: userID, Valid: true},
        Date:        sql.NullTime{Time: date, Valid: true},
        Time:        sql.NullTime{Time: timeNow, Valid: true},
    })
}

func (r *TodoRepository) GetTodosByUserID(ctx context.Context, userID string) ([]domain.Todo, error) {
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    var result []domain.Todo
    for _, todo := range todos {
        result = append(result, domain.Todo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
        })
    }
    return result, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) error {
    return r.querier.UpdateTodo(ctx, db.UpdateTodoParams{
        ID:          id,
        Task:        task,
        Description: sql.NullString{String: description, Valid: true},
        Done:        done,
        Important:   important,
        UserID:      sql.NullString{String: userID, Valid: true},
    })
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id, userID string) error {
    return r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
}