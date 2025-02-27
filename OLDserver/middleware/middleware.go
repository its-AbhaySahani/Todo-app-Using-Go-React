package middleware

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "net/http"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/gorilla/mux"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Get all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    todos, err := models.GetTodos(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(todos)
}

// Create a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var todo models.Todo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    createdTodo, err := models.CreateTodo(todo.Task, todo.Description, todo.Important, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTodo)
}

// Update an existing todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    var todo models.Todo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    updatedTodo, err := models.UpdateTodo(params["id"], todo.Task, todo.Description, todo.Done, todo.Important, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(updatedTodo)
}

// Delete a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Delete associated routines first
    if err := models.DeleteRoutinesByTaskID(params["id"]); err != nil {
        http.Error(w, "Failed to delete associated routines: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Then delete the todo
    err := models.DeleteTodo(params["id"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
// Undo a todo
func UndoTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    undoneTodo, err := models.UndoTodo(params["id"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(undoneTodo)
}

// Share a todo with another user
func ShareTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        TaskID   string `json:"taskId"`
        Username string `json:"username"`
    }
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    user, err := models.GetUserByUsername(request.Username)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    err = models.ShareTodoWithUser(request.TaskID, user.ID, userID)
    if err != nil {
        http.Error(w, "Error sharing task", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Get shared todos
func GetSharedTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    sharedTodos, err := models.GetSharedTodos(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    sharedByMeTodos, err := models.GetSharedByMeTodos(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := map[string]interface{}{
        "received": sharedTodos,
        "shared":   sharedByMeTodos,
    }
    json.NewEncoder(w).Encode(response)
}

// Team-related handlers

// Create a new team
func CreateTeam(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var team models.Team
    _ = json.NewDecoder(r.Body).Decode(&team)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    createdTeam, err := models.CreateTeam(team.Name, team.Password, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTeam)
}

// Join an existing team
func JoinTeam(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        TeamName string `json:"teamName"`
        Password string `json:"password"`
    }
    _ = json.NewDecoder(r.Body).Decode(&request)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    err := models.JoinTeam(request.TeamName, request.Password, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Create a new team todo
func CreateTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var todo models.TeamTodo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    params := mux.Vars(r)
    createdTodo, err := models.CreateTeamTodo(todo.Task, todo.Description, todo.Important, params["teamId"], userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(createdTodo)
}

// Get all teams for the authenticated user
func GetTeams(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    teams, err := models.GetTeams(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(teams)
}

// Get all team todos
func GetTeamTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    todos, err := models.GetTeamTodos(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(todos)
}

// Update a team todo
func UpdateTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    var todo models.TeamTodo
    _ = json.NewDecoder(r.Body).Decode(&todo)
    updatedTodo, err := models.UpdateTeamTodo(params["id"], todo.Task, todo.Description, todo.Done, todo.Important, todo.TeamID, todo.AssignedTo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(updatedTodo)
}

// Delete a team todo
func DeleteTeamTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    err := models.DeleteTeamTodo(params["id"], params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Remove a team member
func RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    err := models.RemoveTeamMember(params["teamId"], params["userId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// Get team details
func GetTeamDetails(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    team, err := models.GetTeamByID(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    todos, err := models.GetTeamTodos(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := map[string]interface{}{
        "team":  team,
        "tasks": todos,
    }
    json.NewEncoder(w).Encode(response)
}

// Get team members
func GetTeamMembers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    members, err := models.GetTeamMembers(params["teamId"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(members)
}

// Add a team member
func AddTeamMember(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var request struct {
        Username string `json:"username"`
    }
    _ = json.NewDecoder(r.Body).Decode(&request)
    params := mux.Vars(r)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    err := models.AddTeamMember(params["teamId"], request.Username, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

// CreateOrUpdateRoutines updates the routines for a task
func CreateOrUpdateRoutines(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var request struct {
        TaskID    string   `json:"taskId"`
        Schedules []string `json:"schedules"` // morning, noon, evening, night
        Day       string   `json:"day"`       // day of the week
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Debug logging
    fmt.Printf("CreateOrUpdateRoutines called with taskID=%s, schedules=%v, day=%s\n", 
        request.TaskID, request.Schedules, request.Day)
    
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Use the provided day or default to current day
    dayName := request.Day
    if dayName == "" {
        dayName = strings.ToLower(time.Now().Weekday().String())
    }

    // First, get existing routines for this task and user
    existingRoutines, err := models.GetRoutinesByTaskID(request.TaskID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Filter routines for this user
    var userRoutines []models.Routine
    for _, routine := range existingRoutines {
        if routine.UserID == userID {
            userRoutines = append(userRoutines, routine)
        }
    }
    
    // Track which schedule types we find and update
    scheduleTypesProcessed := make(map[string]bool)
    var createdRoutines []models.Routine
    
    // First, process existing routines - update or deactivate them
    for i := range userRoutines {
        routine := userRoutines[i]
        // If this schedule type is in the request for any day
        if contains(request.Schedules, routine.ScheduleType) {
            // This schedule type is requested, check if day needs updating
            if routine.Day != dayName {
                // Update the day for this routine
                fmt.Printf("Updating day for routine: ID=%s, scheduleType=%s, old day=%s, new day=%s\n",
                    routine.ID, routine.ScheduleType, routine.Day, dayName)
                
                if err := models.UpdateRoutineDay(routine.ID, dayName); err != nil {
                    http.Error(w, "Failed to update routine day: "+err.Error(), http.StatusInternalServerError)
                    return
                }
                // Update our copy as well
                routine.Day = dayName
            }
            
            // Make sure the routine is active
            if !routine.IsActive {
                if err := models.UpdateRoutineStatus(routine.ID, true); err != nil {
                    http.Error(w, "Failed to update routine status: "+err.Error(), http.StatusInternalServerError)
                    return
                }
                routine.IsActive = true
            }
            
            // Mark this schedule type as processed
            scheduleTypesProcessed[routine.ScheduleType] = true
            createdRoutines = append(createdRoutines, routine)
        } else if routine.Day == dayName {
            // This schedule type is not requested for this day, deactivate it
            fmt.Printf("Deactivating routine: ID=%s, scheduleType=%s, day=%s\n",
                routine.ID, routine.ScheduleType, routine.Day)
            
            if err := models.UpdateRoutineStatus(routine.ID, false); err != nil {
                http.Error(w, "Failed to deactivate routine: "+err.Error(), http.StatusInternalServerError)
                return
            }
        }
    }
    
    // Now create any new routines for schedule types not yet processed
    for _, scheduleType := range request.Schedules {
        if !scheduleTypesProcessed[scheduleType] {
            fmt.Printf("Creating new routine: taskID=%s, scheduleType=%s, day=%s\n",
                request.TaskID, scheduleType, dayName)
            
            // Create a new routine for this schedule type
            newRoutine, err := models.CreateRoutine(dayName, scheduleType, request.TaskID, userID)
            if err != nil {
                http.Error(w, "Failed to create routine: "+err.Error(), http.StatusInternalServerError)
                return
            }
            createdRoutines = append(createdRoutines, newRoutine)
        }
    }
    
    json.NewEncoder(w).Encode(createdRoutines)
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// GetTodayRoutines returns all tasks for today's routines by schedule type
func GetTodayRoutines(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    params := mux.Vars(r)
    scheduleType := params["scheduleType"]
    
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    todos, err := models.GetTodayRoutines(scheduleType, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(todos)
}

// GetRoutinesByTask returns all routines for a specific task
func GetRoutinesByTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    params := mux.Vars(r)
    taskID := params["taskId"]
    
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    routines, err := models.GetRoutinesByTaskID(taskID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Filter routines for this user only
    var userRoutines []models.Routine
    for _, routine := range routines {
        if routine.UserID == userID {
            userRoutines = append(userRoutines, routine)
        }
    }
    
    json.NewEncoder(w).Encode(userRoutines)
}


// Add this new handler function
func GetRoutinesByDayAndSchedule(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    params := mux.Vars(r)
    day := params["day"]
    scheduleType := params["scheduleType"]
    
    // Debug logging
    fmt.Printf("GetRoutinesByDayAndSchedule called with day=%s, scheduleType=%s\n", day, scheduleType)
    
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    todos, err := models.GetDailyRoutines(day, scheduleType, userID)
    if err != nil {
        // More debug logging
        fmt.Printf("Error in GetDailyRoutines: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(todos)
}