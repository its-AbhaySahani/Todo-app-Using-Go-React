package todos_repository

import (
    "context"
    "database/sql"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "time"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure TodoRepository implements domain.TodoRepository
var _ domain.TodoRepository = (*TodoRepository)(nil)

type TodoRepository struct {
    querier *db.Queries
    db      *sql.DB
}

// In server/persistent/todos_repository/todos_persistent.go
func (r *TodoRepository) CreateTodo(ctx context.Context, task, description string, done, important bool, userID string, date, todoTime time.Time) (string, error) {
    id := uuid.New().String()
    
    // Ensure the date is valid
    var nullDate sql.NullTime
    if date.IsZero() || date.Year() < 1 || date.Year() > 9999 {
        nullDate = sql.NullTime{Time: time.Now(), Valid: true}
    } else {
        nullDate = sql.NullTime{Time: date, Valid: true}
    }
    
    // Ensure the time is valid
    var nullTime sql.NullTime
    if todoTime.IsZero() || todoTime.Year() < 1 || todoTime.Year() > 9999 {
        now := time.Now()
        validTime := time.Date(2000, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
        nullTime = sql.NullTime{Time: validTime, Valid: true}
    } else {
        nullTime = sql.NullTime{Time: todoTime, Valid: true}
    }
    
    err := r.querier.CreateTodo(ctx, db.CreateTodoParams{
        ID:          id,
        Task:        task,
        Description: sql.NullString{String: description, Valid: true},
        Done:        done,
        Important:   important,
        UserID:      sql.NullString{String: userID, Valid: true},
        Date:        nullDate,
        Time:        nullTime,
    })
    
    if err != nil {
        return "", err
    }
    
    return id, nil
}   

func (r *TodoRepository) GetTodoByID(ctx context.Context, id string) (*domain.Todo, error) {
    // Using a query to get a todo by ID
    var task string
    var description sql.NullString
    var done bool
    var important bool
    var userID sql.NullString
    var date sql.NullString
    var timeValue sql.NullString
    
    err := r.db.QueryRowContext(ctx, "SELECT task, description, done, important, user_id, CAST(date AS CHAR), CAST(time AS CHAR) FROM todos WHERE id = ?", id).
        Scan(&task, &description, &done, &important, &userID, &date, &timeValue)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("todo not found")
        }
        return nil, err
    }
    
    // Parse date
    var dateTime time.Time
    if date.Valid {
        parsedDate, err := time.Parse("2006-01-02", date.String)
        if err == nil {
            dateTime = parsedDate
        }
    }
    
    // Parse time
    var timeVal time.Time
    if timeValue.Valid {
        parsedTime, err := time.Parse("15:04:05", timeValue.String)
        if err == nil {
            // Use a valid year for the time
            hour, min, sec := parsedTime.Clock()
            timeVal = time.Date(2000, 1, 1, hour, min, sec, 0, time.UTC)
        }
    }
    
    return &domain.Todo{
        ID:          id,
        Task:        task,
        Description: description.String,
        Done:        done,
        Important:   important,
        UserID:      userID.String,
        Date:        dateTime,
        Time:        timeVal,
    }, nil
}

func (r *TodoRepository) GetTodosByUserID(ctx context.Context, userID string) ([]domain.Todo, error) {
    todos, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Log the types to understand what we're dealing with
    if len(todos) > 0 {
        fmt.Printf("Date type: %T\n", todos[0].Date)
        fmt.Printf("Time type: %T\n", todos[0].Time)
    }
    
    // Convert db.Todo to domain.Todo
    domainTodos := make([]domain.Todo, len(todos))
    for i, todo := range todos {
        domainTodos[i] = domain.Todo{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            // Skip date/time for now to see if code works otherwise
        }
        
        // Add a type switch to handle different date/time types
        switch d := todo.Date.(type) {
        case time.Time:
            domainTodos[i].Date = d
        case string:
            parsed, _ := time.Parse("2006-01-02", d)
            domainTodos[i].Date = parsed
        }
        
        switch t := todo.Time.(type) {
        case time.Time:
            domainTodos[i].Time = t
        case string:
            parsed, _ := time.Parse("15:04:05", t)
            domainTodos[i].Time = parsed
        }
    }
    
    return domainTodos, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id, task, description string, done, important bool, userID string) (bool, error) {
    // Use your existing DTO and converter
    req := &dto.UpdateTodoRequest{
        ID:          id,
        Task:        task,
        Description: description,
        Done:        done,
        Important:   important,
        UserID:      userID,
    }
    
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id, userID string) (bool, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TodoRepository) UndoTodo(ctx context.Context, id, userID string) (bool, error) {
    err := r.querier.UndoTodo(ctx, db.UndoTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

// Original methods for backward compatibility - you can keep these if you want
func (r *TodoRepository) CreateTodoWithDTO(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTodoDomainRequestToPersistentRequest()
    err := r.querier.CreateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TodoRepository) GetTodosByUserIDWithDTO(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    todosRows, err := r.querier.GetTodosByUserID(ctx, sql.NullString{String: userID, Valid: true})
    if err != nil {
        return nil, err
    }
    
    // Convert []db.GetTodosByUserIDRow to []dto.TodoResponse
    var todoResponses []dto.TodoResponse
    for _, todo := range todosRows {
        // Handle date and time based on their types
        var date time.Time
        var timeVal time.Time
        
        // Handle date conversion
        switch d := todo.Date.(type) {
        case time.Time:
            date = d
        case string:
            parsed, _ := time.Parse("2006-01-02", d)
            date = parsed
        }
        
        // Handle time conversion
        switch t := todo.Time.(type) {
        case time.Time:
            timeVal = t
        case string:
            parsed, _ := time.Parse("15:04:05", t)
            timeVal = parsed
        }
        
        todoResponses = append(todoResponses, dto.TodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description.String,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID.String,
            Date:        date,
            Time:        timeVal,
        })
    }
    
    return &dto.TodosResponse{Todos: todoResponses}, nil
}

func (r *TodoRepository) UpdateTodoWithDTO(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertUpdateTodoDomainRequestToPersistentRequest()
    err := r.querier.UpdateTodo(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) DeleteTodoWithDTO(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.DeleteTodo(ctx, db.DeleteTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TodoRepository) UndoTodoWithDTO(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.UndoTodo(ctx, db.UndoTodoParams{
        ID:     id,
        UserID: sql.NullString{String: userID, Valid: true},
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}