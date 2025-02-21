package team_todos_repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type TeamTodoRepository struct {
    querier *db.Queries
}

func (r *TeamTodoRepository) CreateTeamTodo(ctx context.Context, task, description string, done, important bool, teamID, assignedTo string) error {
    id := uuid.New().String()
    date := time.Now()
    timeNow := time.Now()
    return r.querier.CreateTeamTodo(ctx, db.CreateTeamTodoParams{
        ID:          id,
        Task:        task,
        Description: sql.NullString{String: description, Valid: true},
        Done:        done,
        Important:   sql.NullBool{Bool: important, Valid: true},
        TeamID:      teamID,
        AssignedTo:  sql.NullString{String: assignedTo, Valid: true},
        Date:        sql.NullTime{Time: date, Valid: true},
        Time:        sql.NullTime{Time: timeNow, Valid: true},
    })
}

func (r *TeamTodoRepository) GetTeamTodos(ctx context.Context, teamID string) ([]domain.TeamTodo, error) {
    todos, err := r.querier.GetTeamTodos(ctx, teamID)
    if err != nil {
        return nil, err
    }
    var result []domain.TeamTodo
    for _, todo := range todos {
        result = append(result, domain.TeamTodo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important.Bool,
            TeamID:      todo.TeamID,
            AssignedTo:  todo.AssignedTo.String,
            Date:        todo.Date.Time,
            Time:        todo.Time.Time,
        })
    }
    return result, nil
}

func (r *TeamTodoRepository) UpdateTeamTodo(ctx context.Context, id, task, description string, done, important bool, teamID, assignedTo string) error {
    return r.querier.UpdateTeamTodo(ctx, db.UpdateTeamTodoParams{
        ID:          id,
        Task:        task,
        Description: sql.NullString{String: description, Valid: true},
        Done:        done,
        Important:   sql.NullBool{Bool: important, Valid: true},
        TeamID:      teamID,
        AssignedTo:  sql.NullString{String: assignedTo, Valid: true},
    })
}

func (r *TeamTodoRepository) DeleteTeamTodo(ctx context.Context, id, teamID string) error {
    return r.querier.DeleteTeamTodo(ctx, db.DeleteTeamTodoParams{
        ID:     id,
        TeamID: teamID,
    })
}