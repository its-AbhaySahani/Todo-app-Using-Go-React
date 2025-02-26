package dto

type CreateResponse struct {
    ID string `json:"id"`
}

type SuccessResponse struct {
    Success bool `json:"success"`
}

type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

type TodoResponse struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
}

type TodosResponse struct {
    Todos []TodoResponse `json:"todos"`
}

type TeamResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type TeamsResponse struct {
    Teams []TeamResponse `json:"teams"`
}

type TeamTodoResponse struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    TeamID      string `json:"team_id"`
    AssignedTo  string `json:"assigned_to"`
}

type TeamTodosResponse struct {
    TeamTodos []TeamTodoResponse `json:"team_todos"`
}

type TeamMembersResponse struct {
    Members []UserResponse `json:"members"`
}

type SharedTodosResponse struct {
    SharedTodos []TodoResponse `json:"shared_todos"`
}