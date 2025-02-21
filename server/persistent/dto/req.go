package dto

// Shared Todos Converters
func (req *CreateSharedTodoRequest) ToParams() *db.CreateSharedTodoParams {
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

// Team Members Converters
func (req *AddTeamMemberRequest) ToParams() *db.AddTeamMemberParams {
    return &db.AddTeamMemberParams{
        TeamID:  req.TeamID,
        UserID:  req.UserID,
        IsAdmin: sql.NullBool{Bool: req.IsAdmin, Valid: true},
    }
}

func (req *RemoveTeamMemberRequest) ToParams() *db.RemoveTeamMemberParams {
    return &db.RemoveTeamMemberParams{
        TeamID: req.TeamID,
        UserID: req.UserID,
    }
}

// Team Todos Converters
func (req *CreateTeamTodoRequest) ToParams() *db.CreateTeamTodoParams {
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

func (req *UpdateTeamTodoRequest) ToParams() *db.UpdateTeamTodoParams {
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

// Teams Converters
func (req *CreateTeamRequest) ToParams() *db.CreateTeamParams {
    return &db.CreateTeamParams{
        ID:       uuid.New().String(),
        Name:     req.Name,
        Password: req.Password,
        AdminID:  req.AdminID,
    }
}

// Todos Converters
func (req *CreateTodoRequest) ToParams() *db.CreateTodoParams {
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

func (req *UpdateTodoRequest) ToParams() *db.UpdateTodoParams {
    return &db.UpdateTodoParams{
        ID:          req.ID,
        Task:        req.Task,
        Description: sql.NullString{String: req.Description, Valid: true},
        Done:        req.Done,
        Important:   req.Important,
        UserID:      sql.NullString{String: req.UserID, Valid: true},
    }
}

// Users Converters
func (req *CreateUserRequest) ToParams() *db.CreateUserParams {
    return &db.CreateUserParams{
        ID:       uuid.New().String(),
        Username: req.Username,
        Password: req.Password,
    }
}