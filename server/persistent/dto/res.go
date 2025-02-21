package dto

import (
    "time"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

// Shared Todos Responses
type SharedTodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    UserID      string    `json:"user_id"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
    SharedBy    string    `json:"shared_by"`
}

type SharedTodosResponse struct {
    Received []SharedTodoResponse `json:"received"`
    Shared   []SharedTodoResponse `json:"shared"`
}

// Team Members Responses
type TeamMemberResponse struct {
    TeamID  string `json:"team_id"`
    UserID  string `json:"user_id"`
    IsAdmin bool   `json:"is_admin"`
}

type TeamMembersResponse struct {
    Members []TeamMemberResponse `json:"members"`
}

// Team Todos Responses
type TeamTodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    TeamID      string    `json:"team_id"`
    AssignedTo  string    `json:"assigned_to"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
}

type TeamTodosResponse struct {
    Todos []TeamTodoResponse `json:"todos"`
}

// Teams Responses
type TeamResponse struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Password string `json:"password,omitempty"`
    AdminID  string `json:"admin_id"`
}

type TeamsResponse struct {
    Teams []TeamResponse `json:"teams"`
}

// Todos Responses
type TodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    UserID      string    `json:"user_id"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
}

type TodosResponse struct {
    Todos []TodoResponse `json:"todos"`
}

// Users Responses
type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password,omitempty"`
}

// Success Responses
type SuccessResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
}

type CreateResponse struct {
    ID string `json:"id"`
}

// Converters
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

func NewUserResponse(user *db.User) *UserResponse {
    return &UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Password: user.Password,
    }
}