package dto

import "time"

type CreateResponse struct {
    ID string `json:"id"`
}

type SuccessResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
}

type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password,omitempty"`
}

type TodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    UserID      string    `json:"userId"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
}

type TodosResponse struct {
    Todos []TodoResponse `json:"todos"`
}

type TeamResponse struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Password string `json:"password,omitempty"`
    AdminID  string `json:"adminId"`
}

type TeamsResponse struct {
    Teams []TeamResponse `json:"teams"`
}

type TeamTodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    TeamID      string    `json:"teamId"`
    AssignedTo  string    `json:"assignedTo"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
}

type TeamTodosResponse struct {
    Todos []TeamTodoResponse `json:"todos"`
}

type TeamMemberResponse struct {
    TeamID  string `json:"teamId"`
    UserID  string `json:"userId"`
    IsAdmin bool   `json:"isAdmin"`
}

type TeamMembersResponse struct {
    Members []TeamMemberResponse `json:"members"`
}

type SharedTodoResponse struct {
    ID          string    `json:"id"`
    Task        string    `json:"task"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Important   bool      `json:"important"`
    UserID      string    `json:"userId"`
    Date        time.Time `json:"date"`
    Time        time.Time `json:"time"`
    SharedBy    string    `json:"sharedBy"`
}

type SharedTodosResponse struct {
    Received []SharedTodoResponse `json:"received"`
    Shared   []SharedTodoResponse `json:"shared"`
}

// Routine Responses
type RoutineResponse struct {
    ID           string    `json:"id"`
    Day          string    `json:"day"`
    ScheduleType string    `json:"scheduleType"`
    TaskID       string    `json:"taskId"`
    UserID       string    `json:"userId"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    IsActive     bool      `json:"isActive"`
}

type RoutinesResponse struct {
    Routines []RoutineResponse `json:"routines"`
}