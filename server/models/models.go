package models

import (
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "golang.org/x/crypto/bcrypt"
    "time"
)

type Todo struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"user_id"`
    Date        string `json:"date"`
    Time        string `json:"time"`
}

type SharedTodo struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    UserID      string `json:"user_id"`
    Date        string `json:"date"`
    Time        string `json:"time"`
    SharedBy    string `json:"shared_by"`
}

type Team struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Password string `json:"password"`
    AdminID  string `json:"admin_id"`
}

type TeamMember struct {
    TeamID  string `json:"team_id"`
    UserID  string `json:"user_id"`
    IsAdmin bool   `json:"is_admin"`
}

type TeamTodo struct {
    ID          string `json:"id"`
    Task        string `json:"task"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    Important   bool   `json:"important"`
    TeamID      string `json:"team_id"`
    AssignedTo  string `json:"assigned_to"`
    Date        string `json:"date"`
    Time        string `json:"time"`
}

type TeamMemberDetails struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
}


func GetTodos(userID string) ([]Todo, error) {
    rows, err := database.DB.Query("SELECT id, task, description, done, important, date, time FROM todos WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Done, &todo.Important, &todo.Date, &todo.Time); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }
    return todos, nil
}

func CreateTodo(task, description string, important bool, userID string) (Todo, error) {
    id := uuid.New().String()
    currentDate := time.Now().Format("2006-01-02")
    currentTime := time.Now().Format("15:04:05")
    _, err := database.DB.Exec("INSERT INTO todos (id, task, description, done, important, user_id, date, time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", id, task, description, false, important, userID, currentDate, currentTime)
    if err != nil {
        return Todo{}, err
    }
    return Todo{ID: id, Task: task, Description: description, Done: false, Important: important, UserID: userID, Date: currentDate, Time: currentTime}, nil
}

func UpdateTodo(id, task, description string, done, important bool, userID string) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET task = ?, description = ?, done = ?, important = ? WHERE id = ? AND user_id = ?", task, description, done, important, id, userID)
    if err != nil {
        return Todo{}, err
    }
    var todo Todo
    err = database.DB.QueryRow("SELECT id, task, description, done, important, date, time FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Done, &todo.Important, &todo.Date, &todo.Time)
    if err != nil {
        return Todo{}, err
    }
    return todo, nil
}

func DeleteTodo(id, userID string) error {
    _, err := database.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", id, userID)
    return err
}

func UndoTodo(id, userID string) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET done = ? WHERE id = ? AND user_id = ?", false, id, userID)
    if err != nil {
        return Todo{}, err
    }
    var todo Todo
    err = database.DB.QueryRow("SELECT id, task, description, done, important, date, time FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Done, &todo.Important, &todo.Date, &todo.Time)
    if err != nil {
        return Todo{}, err
    }
    return todo, nil
}

func ShareTodoWithUser(taskID, userID, sharedBy string) error {
    _, err := database.DB.Exec("INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by) SELECT id, task, description, done, important, ?, date, time, ? FROM todos WHERE id = ?", userID, sharedBy, taskID)
    return err
}

func GetSharedTodos(userID string) ([]SharedTodo, error) {
    rows, err := database.DB.Query("SELECT id, task, description, done, important, user_id, date, time, shared_by FROM shared_todos WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sharedTodos []SharedTodo
    for rows.Next() {
        var sharedTodo SharedTodo
        if err := rows.Scan(&sharedTodo.ID, &sharedTodo.Task, &sharedTodo.Description, &sharedTodo.Done, &sharedTodo.Important, &sharedTodo.UserID, &sharedTodo.Date, &sharedTodo.Time, &sharedTodo.SharedBy); err != nil {
            return nil, err
        }
        sharedTodos = append(sharedTodos, sharedTodo)
    }
    return sharedTodos, nil
}

func GetSharedByMeTodos(userID string) ([]SharedTodo, error) {
    rows, err := database.DB.Query("SELECT id, task, description, done, important, user_id, date, time, shared_by FROM shared_todos WHERE shared_by = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sharedByMeTodos []SharedTodo
    for rows.Next() {
        var sharedTodo SharedTodo
        if err := rows.Scan(&sharedTodo.ID, &sharedTodo.Task, &sharedTodo.Description, &sharedTodo.Done, &sharedTodo.Important, &sharedTodo.UserID, &sharedTodo.Date, &sharedTodo.Time, &sharedTodo.SharedBy); err != nil {
            return nil, err
        }
        sharedByMeTodos = append(sharedByMeTodos, sharedTodo)
    }
    return sharedByMeTodos, nil
}

// Team-related functions

func CreateTeam(name, password, adminID string) (Team, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return Team{}, err
    }

    id := uuid.New().String()
    _, err = database.DB.Exec("INSERT INTO teams (id, name, password, admin_id) VALUES (?, ?, ?, ?)", id, name, string(hashedPassword), adminID)
    if err != nil {
        return Team{}, err
    }

    _, err = database.DB.Exec("INSERT INTO team_members (team_id, user_id, is_admin) VALUES (?, ?, ?)", id, adminID, true)
    if err != nil {
        return Team{}, err
    }

    return Team{ID: id, Name: name, Password: string(hashedPassword), AdminID: adminID}, nil
}

func JoinTeam(teamID, password, userID string) error {
    var hashedPassword string
    err := database.DB.QueryRow("SELECT password FROM teams WHERE id = ?", teamID).Scan(&hashedPassword)
    if err != nil {
        return err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return err
    }

    _, err = database.DB.Exec("INSERT INTO team_members (team_id, user_id) VALUES (?, ?)", teamID, userID)
    return err
}

func CreateTeamTodo(task, description string, important bool, teamID, assignedTo string) (TeamTodo, error) {
    id := uuid.New().String()
    currentDate := time.Now().Format("2006-01-02")
    currentTime := time.Now().Format("15:04:05")
    _, err := database.DB.Exec("INSERT INTO team_todos (id, task, description, done, important, team_id, assigned_to, date, time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", id, task, description, false, important, teamID, assignedTo, currentDate, currentTime)
    if err != nil {
        return TeamTodo{}, err
    }
    return TeamTodo{ID: id, Task: task, Description: description, Done: false, Important: important, TeamID: teamID, AssignedTo: assignedTo, Date: currentDate, Time: currentTime}, nil
}

func GetTeamTodos(teamID string) ([]TeamTodo, error) {
    rows, err := database.DB.Query("SELECT id, task, description, done, important, team_id, assigned_to, date, time FROM team_todos WHERE team_id = ?", teamID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var todos []TeamTodo
    for rows.Next() {
        var todo TeamTodo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Done, &todo.Important, &todo.TeamID, &todo.AssignedTo, &todo.Date, &todo.Time); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }
    return todos, nil
}

func GetTeams(userID string) ([]Team, error) {
    rows, err := database.DB.Query("SELECT t.id, t.name, t.password, t.admin_id FROM teams t JOIN team_members tm ON t.id = tm.team_id WHERE tm.user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teams []Team
    for rows.Next() {
        var team Team
        if err := rows.Scan(&team.ID, &team.Name, &team.Password, &team.AdminID); err != nil {
            return nil, err
        }
        teams = append(teams, team)
    }
    return teams, nil
}


func UpdateTeamTodo(id, task, description string, done, important bool, teamID, assignedTo string) (TeamTodo, error) {
    _, err := database.DB.Exec("UPDATE team_todos SET task = ?, description = ?, done = ?, important = ?, assigned_to = ? WHERE id = ? AND team_id = ?", task, description, done, important, assignedTo, id, teamID)
    if err != nil {
        return TeamTodo{}, err
    }
    var todo TeamTodo
    err = database.DB.QueryRow("SELECT id, task, description, done, important, team_id, assigned_to, date, time FROM team_todos WHERE id = ? AND team_id = ?", id, teamID).Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Done, &todo.Important, &todo.TeamID, &todo.AssignedTo, &todo.Date, &todo.Time)
    if err != nil {
        return TeamTodo{}, err
    }
    return todo, nil
}

func DeleteTeamTodo(id, teamID string) error {
    _, err := database.DB.Exec("DELETE FROM team_todos WHERE id = ? AND team_id = ?", id, teamID)
    return err
}

func RemoveTeamMember(teamID, userID string) error {
    _, err := database.DB.Exec("DELETE FROM team_members WHERE team_id = ? AND user_id = ?", teamID, userID)
    return err
}

func GetTeamByID(teamID string) (Team, error) {
    var team Team
    err := database.DB.QueryRow("SELECT id, name, password, admin_id FROM teams WHERE id = ?", teamID).Scan(&team.ID, &team.Name, &team.Password, &team.AdminID)
    if err != nil {
        return Team{}, err
    }
    return team, nil
}

func GetTeamMembers(teamID string) ([]TeamMemberDetails, error) {
    rows, err := database.DB.Query("SELECT u.id, u.username, tm.is_admin FROM users u JOIN team_members tm ON u.id = tm.user_id WHERE tm.team_id = ?", teamID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var members []TeamMemberDetails
    for rows.Next() {
        var member TeamMemberDetails
        if err := rows.Scan(&member.ID, &member.Username, &member.IsAdmin); err != nil {
            return nil, err
        }
        members = append(members, member)
    }
    return members, nil
}