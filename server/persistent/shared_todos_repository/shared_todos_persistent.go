package shared_todos_repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure SharedTodoRepository implements domain.SharedTodoRepository
var _ domain.SharedTodoRepository = (*SharedTodoRepository)(nil)

type SharedTodoRepository struct {
    querier *db.Queries
    db      *sql.DB  // Add this field
}
// Implement domain.SharedTodoRepository interface methods
func (r *SharedTodoRepository) CreateSharedTodo(ctx context.Context, task, description string, done, important bool, userID, sharedBy string) (string, error) {
    id := uuid.New().String()
    
    // Use current date and time if none provided
    nullDate := sql.NullTime{Time: time.Now(), Valid: true}
    nullTime := sql.NullTime{Time: time.Now(), Valid: true}
    
    err := r.querier.CreateSharedTodo(ctx, db.CreateSharedTodoParams{
        ID:          id,
        Task:        sql.NullString{String: task, Valid: true},
        Description: sql.NullString{String: description, Valid: true},
        Done:        sql.NullBool{Bool: done, Valid: true},
        Important:   sql.NullBool{Bool: important, Valid: true},
        UserID:      sql.NullString{String: userID, Valid: true},
        Date:        nullDate,
        Time:        nullTime,
        SharedBy:    sql.NullString{String: sharedBy, Valid: true},
    })
    
    if err != nil {
        return "", err
    }
    
    return id, nil
}

func (r *SharedTodoRepository) GetSharedTodos(ctx context.Context, userID string) ([]domain.SharedTodo, error) {
    // Cast date and time to CHAR to get them as strings
    rows, err := r.db.QueryContext(ctx, 
        "SELECT id, task, description, done, important, user_id, CAST(date AS CHAR), CAST(time AS CHAR), shared_by FROM shared_todos WHERE user_id = ?", 
        userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var todos []domain.SharedTodo
    for rows.Next() {
        var todo domain.SharedTodo
        var description, dateStr, timeStr sql.NullString
        
        if err := rows.Scan(
            &todo.ID,
            &todo.Task,
            &description,
            &todo.Done,
            &todo.Important,
            &todo.UserID,
            &dateStr,
            &timeStr,
            &todo.SharedBy,
        ); err != nil {
            return nil, err
        }
        
        todo.Description = description.String
        
        // Parse date
        if dateStr.Valid && dateStr.String != "" {
            parsed, err := time.Parse("2006-01-02", dateStr.String)
            if err == nil {
                todo.Date = parsed
            } else {
                todo.Date = time.Now() // Default to today if parsing fails
            }
        } else {
            todo.Date = time.Now() // Default to today if no date
        }
        
        // Parse time
        if timeStr.Valid && timeStr.String != "" {
            parsed, err := time.Parse("15:04:05", timeStr.String)
            if err == nil {
                // Use a valid year for the time
                hour, min, sec := parsed.Clock()
                todo.Time = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            } else {
                // Default to current time if parsing fails
                now := time.Now()
                todo.Time = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            }
        } else {
            // Default to current time if no time
            now := time.Now()
            todo.Time = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        }
        
        todos = append(todos, todo)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return todos, nil
}


func (r *SharedTodoRepository) GetSharedByMeTodos(ctx context.Context, sharedBy string) ([]domain.SharedTodo, error) {
    // Cast date and time to CHAR to get them as strings
    rows, err := r.db.QueryContext(ctx, 
        "SELECT id, task, description, done, important, user_id, CAST(date AS CHAR), CAST(time AS CHAR), shared_by FROM shared_todos WHERE shared_by = ?", 
        sharedBy)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var todos []domain.SharedTodo
    for rows.Next() {
        var todo domain.SharedTodo
        var description, dateStr, timeStr sql.NullString
        
        if err := rows.Scan(
            &todo.ID,
            &todo.Task,
            &description,
            &todo.Done,
            &todo.Important,
            &todo.UserID,
            &dateStr,
            &timeStr,
            &todo.SharedBy,
        ); err != nil {
            return nil, err
        }
        
        todo.Description = description.String
        
        // Parse date
        if dateStr.Valid && dateStr.String != "" {
            parsed, err := time.Parse("2006-01-02", dateStr.String)
            if err == nil {
                todo.Date = parsed
            } else {
                todo.Date = time.Now() // Default to today if parsing fails
            }
        } else {
            todo.Date = time.Now() // Default to today if no date
        }
        
        // Parse time
        if timeStr.Valid && timeStr.String != "" {
            parsed, err := time.Parse("15:04:05", timeStr.String)
            if err == nil {
                // Use a valid year for the time
                hour, min, sec := parsed.Clock()
                todo.Time = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            } else {
                // Default to current time if parsing fails
                now := time.Now()
                todo.Time = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            }
        } else {
            // Default to current time if no time
            now := time.Now()
            todo.Time = time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        }
        
        todos = append(todos, todo)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return todos, nil
}

// Original methods for backward compatibility
func (r *SharedTodoRepository) CreateSharedTodoWithDTO(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateSharedTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateSharedTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *SharedTodoRepository) GetSharedTodosWithDTO(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedTodos(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}

func (r *SharedTodoRepository) GetSharedByMeTodosWithDTO(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    todos, err := r.querier.GetSharedByMeTodos(ctx, sql.NullString{String: sharedBy, Valid: true})
    if err != nil {
        return nil, err
    }
    return dto.NewSharedTodosResponse(todos), nil
}

// ShareTodo shares a todo with another user
func (r *SharedTodoRepository) ShareTodo(ctx context.Context, todoID string, recipientUserID string, sharedBy string) error {
    // First get the original todo
    var task string
    var description sql.NullString
    var done bool
    var important bool
    var date sql.NullString
    var timeValue sql.NullString
    
    err := r.db.QueryRowContext(ctx, 
        "SELECT task, description, done, important, CAST(date AS CHAR), CAST(time AS CHAR) FROM todos WHERE id = ?", 
        todoID).Scan(&task, &description, &done, &important, &date, &timeValue)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("todo not found")
        }
        return err
    }
    
    // Create a new shared todo
    id := uuid.New().String()
    
    // Prepare date and time values
    var nullDate sql.NullTime
    var nullTime sql.NullTime
    
    // Handle date
    if date.Valid && date.String != "" {
        parsedDate, err := time.Parse("2006-01-02", date.String)
        if err == nil {
            nullDate = sql.NullTime{Time: parsedDate, Valid: true}
        } else {
            nullDate = sql.NullTime{Time: time.Now(), Valid: true}
        }
    } else {
        nullDate = sql.NullTime{Time: time.Now(), Valid: true}
    }
    
    // Handle time
    if timeValue.Valid && timeValue.String != "" {
        parsedTime, err := time.Parse("15:04:05", timeValue.String)
        if err == nil {
            hour, min, sec := parsedTime.Clock()
            validTime := time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
            nullTime = sql.NullTime{Time: validTime, Valid: true}
        } else {
            now := time.Now()
            validTime := time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
            nullTime = sql.NullTime{Time: validTime, Valid: true}
        }
    } else {
        now := time.Now()
        validTime := time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        nullTime = sql.NullTime{Time: validTime, Valid: true}
    }
    
    // Insert the shared todo
    err = r.querier.CreateSharedTodo(ctx, db.CreateSharedTodoParams{
        ID:          id,
        Task:        sql.NullString{String: task, Valid: true},
        Description: description,
        Done:        sql.NullBool{Bool: done, Valid: true},
        Important:   sql.NullBool{Bool: important, Valid: true},
        UserID:      sql.NullString{String: recipientUserID, Valid: true},
        Date:        nullDate,
        Time:        nullTime,
        SharedBy:    sql.NullString{String: sharedBy, Valid: true},
    })
    
    return err
}

// IsSharedWithUser checks if a todo is already shared with a user
func (r *SharedTodoRepository) IsSharedWithUser(ctx context.Context, todoID string, userID string) (bool, error) {
    var count int
    err := r.db.QueryRowContext(ctx, 
        "SELECT COUNT(*) FROM shared_todos WHERE task IN (SELECT task FROM todos WHERE id = ?) AND user_id = ?", 
        todoID, userID).Scan(&count)
    
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}