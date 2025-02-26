package dto

type CreateUserRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

type CreateTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Important   bool   `json:"important"`
}

type UpdateTodoRequest struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
}

type ShareTodoRequest struct {
    TaskID   string `json:"task_id"`
    Username string `json:"username"`
}

type CreateTeamRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type JoinTeamRequest struct {
    TeamName string `json:"team_name"`
    Password string `json:"password"`
}

type CreateTeamTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Important   bool   `json:"important"`
    TeamID      string `json:"team_id"`
    AssignedTo  string `json:"assigned_to"`
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

type AddTeamMemberRequest struct {
    TeamID string `json:"team_id"`
    UserID string `json:"user_id"`
}