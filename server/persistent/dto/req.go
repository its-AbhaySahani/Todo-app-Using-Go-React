package dto

import (
    "database/sql"
    "time"
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

// Shared Todos
type CreateSharedTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"user_id"`
    SharedBy    string `json:"shared_by"`
}

func (req *CreateSharedTodoRequest) ConvertCreateSharedTodoDomainRequestToPersistentRequest() *db.CreateSharedTodoParams {
    return &db.CreateSharedTodoParams{
        ID:          uuid.New().String(),
        Task:        sql.NullString{String: req.Task, Valid: true},
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        sql.NullBool{Bool: req.Done, Valid: true},
        Important:   sql.NullBool{Bool: req.Important, Valid: true},
        UserID:      sql.NullString{String: req.UserID, Valid: true},
        SharedBy:    sql.NullString{String: req.SharedBy, Valid: true},
        Date:        sql.NullTime{Time: time.Now(), Valid: true},
        Time:        sql.NullTime{Time: time.Now(), Valid: true},
    }
}

// Team Members
type AddTeamMemberRequest struct {
    TeamID  string `json:"team_id"`
    UserID  string `json:"user_id"`
    IsAdmin bool   `json:"is_admin"`
}

func (req *AddTeamMemberRequest) ConvertAddTeamMemberDomainRequestToPersistentRequest() *db.AddTeamMemberParams {
    return &db.AddTeamMemberParams{
        TeamID:  req.TeamID,
        UserID:  req.UserID,
        IsAdmin: sql.NullBool{Bool: req.IsAdmin, Valid: true},
    }
}


// Team Todos
type CreateTeamTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    TeamID      string `json:"team_id"`
    AssignedTo  string `json:"assigned_to"`
}

func (req *CreateTeamTodoRequest) ConvertCreateTeamTodoDomainRequestToPersistentRequest() *db.CreateTeamTodoParams {
    return &db.CreateTeamTodoParams{
        ID:          uuid.New().String(),
        Task:        req.Task,
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        req.Done,
        Important:   sql.NullBool{Bool: req.Important, Valid: true},
        TeamID:      req.TeamID,
        AssignedTo:  sql.NullString{String: req.AssignedTo, Valid: true},
        Date:        sql.NullTime{Time: time.Now(), Valid: true},
        Time:        sql.NullTime{Time: time.Now(), Valid: true},
    }
}

type UpdateTeamTodoRequest struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    TeamID      string `json:"team_id"`
    AssignedTo  string `json:"assigned_to"`
}

func (req *UpdateTeamTodoRequest) ConvertUpdateTeamTodoDomainRequestToPersistentRequest() *db.UpdateTeamTodoParams {
    return &db.UpdateTeamTodoParams{
        ID:          req.ID,
        Task:        req.Task,
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        req.Done,
        Important:   sql.NullBool{Bool: req.Important, Valid: true},
        TeamID:      req.TeamID,
        AssignedTo:  sql.NullString{String: req.AssignedTo, Valid: true},
    }
}

// Teams
type CreateTeamRequest struct {
    Name     string `json:"name"`
    Password string `json:"password"`
    AdminID  string `json:"admin_id"`
}

func (req *CreateTeamRequest) ConvertCreateTeamDomainRequestToPersistentRequest() *db.CreateTeamParams {
    return &db.CreateTeamParams{
        ID:       uuid.New().String(),
        Name:     req.Name,
        Password: req.Password,
        AdminID:  req.AdminID,
    }
}

// Todos
type CreateTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"user_id"`
}

func (req *CreateTodoRequest) ConvertCreateTodoDomainRequestToPersistentRequest() *db.CreateTodoParams {
    return &db.CreateTodoParams{
        ID:          uuid.New().String(),
        Task:        req.Task,
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        req.Done,
        Important:   req.Important,
        UserID:      sql.NullString{String: req.UserID, Valid: true},
        Date:        sql.NullTime{Time: time.Now(), Valid: true},
        Time:        sql.NullTime{Time: time.Now(), Valid: true},
    }
}

type UpdateTodoRequest struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"user_id"`
}

func (req *UpdateTodoRequest) ConvertUpdateTodoDomainRequestToPersistentRequest() *db.UpdateTodoParams {
    return &db.UpdateTodoParams{
        ID:          req.ID,
        Task:        req.Task,
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        req.Done,
        Important:   req.Important,
        UserID:      sql.NullString{String: req.UserID, Valid: true},
    }
}

// Users
type CreateUserRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (req *CreateUserRequest) ConvertCreateUserDomainRequestToPersistentRequest() *db.CreateUserParams {
    return &db.CreateUserParams{
        ID:       uuid.New().String(),
        Username: req.Username,
        Password: req.Password,
    }
}


