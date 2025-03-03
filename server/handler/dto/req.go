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
    UserID      string `json:"userId,omitempty"` // Will be set from context
    Done        bool   `json:"done"`
}

type UpdateTodoRequest struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"userId,omitempty"` // Will be set from context
}

type ShareTodoRequest struct {
    TaskID   string `json:"taskId"`
    Username string `json:"username"`
    SharedBy string `json:"sharedBy,omitempty"` // Will be set from context
}

type CreateTeamRequest struct {
    Name     string `json:"name"`
    Password string `json:"password"`
    AdminID  string `json:"adminId,omitempty"` // Will be set from context
}

type JoinTeamRequest struct {
    TeamName string `json:"teamName"`
    Password string `json:"password"`
    UserID   string `json:"userId,omitempty"` // Will be set from context
}

type CreateTeamTodoRequest struct {
    Task        string `json:"task"`
    Description string `json:"description"`
    Important   bool   `json:"important"`
    TeamID      string `json:"teamId,omitempty"` // Will be set from URL params
    AssignedTo  string `json:"assignedTo,omitempty"` // Will be set from context
}

type UpdateTeamTodoRequest struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    TeamID      string `json:"teamId,omitempty"` // Will be set from URL params
    AssignedTo  string `json:"assignedTo,omitempty"` // Will be set from context
}

type AddTeamMemberRequest struct {
    TeamID  string `json:"teamId,omitempty"` // Will be set from URL params
    UserID  string `json:"userId"`
    IsAdmin bool   `json:"isAdmin"`
}

// Routine Requests
type CreateRoutineRequest struct {
    Day          string `json:"day"`
    ScheduleType string `json:"scheduleType"`
    TaskID       string `json:"taskId"`
    UserID       string `json:"userId,omitempty"` // Will be set from context
    IsActive     bool   `json:"isActive"`
}

type UpdateRoutineDayRequest struct {
    ID  string `json:"id"`
    Day string `json:"day"`
}

type UpdateRoutineStatusRequest struct {
    ID       string `json:"id"`
    IsActive bool   `json:"isActive"`
}

type CreateOrUpdateRoutinesRequest struct {
    TaskID    string   `json:"taskId"`
    Schedules []string `json:"schedules"`
    Day       string   `json:"day"`
    UserID    string   `json:"userId,omitempty"` // Will be set from context
}

type GetDailyRoutinesRequest struct {
    Day          string `json:"day"`
    ScheduleType string `json:"scheduleType"`
    UserID       string `json:"userId,omitempty"` // Will be set from context
}

type GetTodayRoutinesRequest struct {
    ScheduleType string `json:"scheduleType"`
    UserID       string `json:"userId,omitempty"` // Will be set from context
}