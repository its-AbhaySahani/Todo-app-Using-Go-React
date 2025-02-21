package dto

// Shared Todos Converters
func NewSharedTodoResponse(todo *db.SharedTodo) *SharedTodoResponse {
    return &SharedTodoResponse{
        ID:          todo.ID,
        Task:        todo.Task.String,
        Description: todo.Description.String,
        Done:        todo.Done.Bool,
        Important:   todo.Important.Bool,
        UserID:      todo.UserID.String,
        Date:        todo.Date.Time,
        Time:        todo.Time.Time,
        SharedBy:    todo.SharedBy.String,
    }
}

func NewSharedTodosResponse(todos []db.SharedTodo) *SharedTodosResponse {
    var response SharedTodosResponse
    for _, todo := range todos {
        response.Received = append(response.Received, *NewSharedTodoResponse(&todo))
    }
    return &response
}

// Team Members Converters
func NewTeamMemberResponse(member *db.TeamMember) *TeamMemberResponse {
    return &TeamMemberResponse{
        TeamID:  member.TeamID,
        UserID:  member.UserID,
        IsAdmin: member.IsAdmin.Bool,
    }
}

func NewTeamMembersResponse(members []db.TeamMember) *TeamMembersResponse {
    var response TeamMembersResponse
    for _, member := range members {
        response.Members = append(response.Members, *NewTeamMemberResponse(&member))
    }
    return &response
}

// Team Todos Converters
func NewTeamTodoResponse(todo *db.TeamTodo) *TeamTodoResponse {
    return &TeamTodoResponse{
        ID:          todo.ID,
        Task:        todo.Task,
        Description: todo.Description.String,
        Done:        todo.Done,
        Important:   todo.Important.Bool,
        TeamID:      todo.TeamID,
        AssignedTo:  todo.AssignedTo.String,
        Date:        todo.Date.Time,
        Time:        todo.Time.Time,
    }
}

func NewTeamTodosResponse(todos []db.TeamTodo) *TeamTodosResponse {
    var response TeamTodosResponse
    for _, todo := range todos {
        response.Todos = append(response.Todos, *NewTeamTodoResponse(&todo))
    }
    return &response
}

// Teams Converters
func NewTeamResponse(team *db.Team) *TeamResponse {
    return &TeamResponse{
        ID:       team.ID,
        Name:     team.Name,
        Password: team.Password,
        AdminID:  team.AdminID,
    }
}

func NewTeamsResponse(teams []db.Team) *TeamsResponse {
    var response TeamsResponse
    for _, team := range teams {
        response.Teams = append(response.Teams, *NewTeamResponse(&team))
    }
    return &response
}

// Todos Converters
func NewTodoResponse(todo *db.Todo) *TodoResponse {
    return &TodoResponse{
        ID:          todo.ID,
        Task:        todo.Task,
        Description: todo.Description.String,
        Done:        todo.Done,
        Important:   todo.Important,
        UserID:      todo.UserID.String,
        Date:        todo.Date.Time,
        Time:        todo.Time.Time,
    }
}

func NewTodosResponse(todos []db.Todo) *TodosResponse {
    var response TodosResponse
    for _, todo := range todos {
        response.Todos = append(response.Todos, *NewTodoResponse(&todo))
    }
    return &response
}

// Users Converters
func NewUserResponse(user *db.User) *UserResponse {
    return &UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Password: user.Password,
    }
}