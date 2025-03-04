package handler

import (
    "database/sql"
    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/api"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/middleware"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/users_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/teams_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_members_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/shared_todos_repository"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/routine_repository"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/users"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/teams"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/team_members"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/team_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/shared_todos"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/services/routines"
)

func SetupRoutes(router *mux.Router, DB *sql.DB) {
    // Initialize repositories
    userRepo := users_repository.NewUserRepository(DB)
    todoRepo := todos_repository.NewTodoRepository(DB)
    teamRepo := teams_repository.NewTeamRepository(DB)
    teamMemberRepo := team_members_repository.NewTeamMemberRepository(DB)
    teamTodoRepo := team_todos_repository.NewTeamTodoRepository(DB)
    sharedTodoRepo := shared_todos_repository.NewSharedTodoRepository(DB)
    routineRepo := routine_repository.NewRoutineRepository(DB)

    // Initialize services
    userService := users.NewUserService(userRepo)
    todoService := todos.NewTodoService(todoRepo)
    teamService := teams.NewTeamService(teamRepo)
    teamMemberService := team_members.NewTeamMemberService(teamMemberRepo)
    teamTodoService := team_todos.NewTeamTodoService(teamTodoRepo)
    sharedTodoService := shared_todos.NewSharedTodoService(sharedTodoRepo, todoRepo, userRepo)
    routineService := routines.NewRoutineService(routineRepo)

    // Setup API v1 routes
    setupV1Routes(router, userService, todoService, teamService, teamMemberService, teamTodoService, sharedTodoService, routineService)
    
    // For backward compatibility, maintain the existing API routes
    // This helps existing clients to continue working while new clients can use v1 API
    setupLegacyRoutes(router, userService, todoService, teamService, teamMemberService, teamTodoService, sharedTodoService, routineService)
}

// setupV1Routes configures the versioned API endpoints
func setupV1Routes(
    router *mux.Router,
    userService *users.UserService,
    todoService *todos.TodoService,
    teamService *teams.TeamService,
    teamMemberService *team_members.TeamMemberService,
    teamTodoService *team_todos.TeamTodoService,
    sharedTodoService *shared_todos.SharedTodoService,
    routineService *routines.RoutineService,
) {
    // API v1
    v1 := router.PathPrefix("/api/v1").Subrouter()
    
    // Public routes
    v1.HandleFunc("/register", api.Register(userService)).Methods("POST")
    v1.HandleFunc("/login", api.Login(userService)).Methods("POST")
    
    // Protected routes
    v1Protected := v1.PathPrefix("").Subrouter()
    v1Protected.Use(middleware.AuthMiddleware)
    
    // Todo routes
    v1Protected.HandleFunc("/todos", api.GetTodos(todoService)).Methods("GET")
    v1Protected.HandleFunc("/todo", api.CreateTodo(todoService)).Methods("POST")
    v1Protected.HandleFunc("/todo/{id}", api.UpdateTodo(todoService)).Methods("PUT")
    v1Protected.HandleFunc("/todo/{id}", api.DeleteTodo(todoService)).Methods("DELETE")
    v1Protected.HandleFunc("/todo/undo/{id}", api.UndoTodo(todoService)).Methods("PUT")
    v1Protected.HandleFunc("/shared", api.GetSharedTodos(sharedTodoService)).Methods("GET")
    
    // Team routes
    v1Protected.HandleFunc("/team", api.CreateTeam(teamService)).Methods("POST")
    v1Protected.HandleFunc("/teams", api.GetTeams(teamService)).Methods("GET")
    v1Protected.HandleFunc("/team/{teamId}/todos", api.GetTeamTodos(teamTodoService)).Methods("GET")
    v1Protected.HandleFunc("/team/{teamId}/todo", api.CreateTeamTodo(teamTodoService)).Methods("POST")
    v1Protected.HandleFunc("/team/{teamId}/todo/{id}", api.UpdateTeamTodo(teamTodoService)).Methods("PUT")
    v1Protected.HandleFunc("/team/{teamId}/todo/{id}", api.DeleteTeamTodo(teamTodoService)).Methods("DELETE")
    v1Protected.HandleFunc("/team/{teamId}/members", api.GetTeamMembers(teamMemberService)).Methods("GET")
    v1Protected.HandleFunc("/team/{teamId}/member", api.AddTeamMember(teamMemberService)).Methods("POST")
    v1Protected.HandleFunc("/team/{teamId}/member/{userId}", api.RemoveTeamMember(teamMemberService)).Methods("DELETE")
// Routine routes
    v1Protected.HandleFunc("/routine", api.CreateOrUpdateRoutines(routineService)).Methods("POST")
    v1Protected.HandleFunc("/routine/task/{taskId}", api.GetRoutinesByTaskID(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/today/{scheduleType}", api.GetTodayRoutines(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/day/{day}/{scheduleType}", api.GetDailyRoutines(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/{id}", api.UpdateRoutineDay(routineService)).Methods("PUT")
    v1Protected.HandleFunc("/routine/{id}/status", api.UpdateRoutineStatus(routineService)).Methods("PUT")
    v1Protected.HandleFunc("/routine/task/{taskId}/delete", api.DeleteRoutinesByTaskID(routineService)).Methods("DELETE")

    // Routine routes
    v1Protected.HandleFunc("/routine", api.CreateOrUpdateRoutines(routineService)).Methods("POST")
    v1Protected.HandleFunc("/routine/task/{taskId}", api.GetRoutinesByTaskID(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/today/{scheduleType}", api.GetTodayRoutines(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/day/{day}/{scheduleType}", api.GetDailyRoutines(routineService)).Methods("GET")
    v1Protected.HandleFunc("/routine/{id}", api.UpdateRoutineDay(routineService)).Methods("PUT")
    v1Protected.HandleFunc("/routine/{id}/status", api.UpdateRoutineStatus(routineService)).Methods("PUT")
    v1Protected.HandleFunc("/routine/task/{taskId}/delete", api.DeleteRoutinesByTaskID(routineService)).Methods("DELETE")

}



// setupLegacyRoutes maintains the original API endpoints for backward compatibility
func setupLegacyRoutes(
    router *mux.Router,
    userService *users.UserService,
    todoService *todos.TodoService,
    teamService *teams.TeamService,
    teamMemberService *team_members.TeamMemberService,
    teamTodoService *team_todos.TeamTodoService,
    sharedTodoService *shared_todos.SharedTodoService,
    routineService *routines.RoutineService,
) {
    // Public routes
    router.HandleFunc("/api/register", api.Register(userService)).Methods("POST")
    router.HandleFunc("/api/login", api.Login(userService)).Methods("POST")
    
    // Protected routes
    apiRouter := router.PathPrefix("/api").Subrouter()
    apiRouter.Use(middleware.AuthMiddleware)
    
    // Todo routes
    apiRouter.HandleFunc("/todos", api.GetTodos(todoService)).Methods("GET")
    apiRouter.HandleFunc("/todo", api.CreateTodo(todoService)).Methods("POST")
    apiRouter.HandleFunc("/todo/{id}", api.UpdateTodo(todoService)).Methods("PUT")
    apiRouter.HandleFunc("/todo/{id}", api.DeleteTodo(todoService)).Methods("DELETE")
    apiRouter.HandleFunc("/todo/undo/{id}", api.UndoTodo(todoService)).Methods("PUT")
    apiRouter.HandleFunc("/shared", api.GetSharedTodos(sharedTodoService)).Methods("GET")
    apiRouter.HandleFunc("/share", api.ShareTodo(sharedTodoService, userService, todoService)).Methods("POST")

    
    // Team routes
    apiRouter.HandleFunc("/team", api.CreateTeam(teamService)).Methods("POST")
    apiRouter.HandleFunc("/teams", api.GetTeams(teamService)).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/todos", api.GetTeamTodos(teamTodoService)).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/todo", api.CreateTeamTodo(teamTodoService)).Methods("POST")
    apiRouter.HandleFunc("/team/{teamId}/todo/{id}", api.UpdateTeamTodo(teamTodoService)).Methods("PUT")
    apiRouter.HandleFunc("/team/{teamId}/todo/{id}", api.DeleteTeamTodo(teamTodoService)).Methods("DELETE")
    apiRouter.HandleFunc("/team/{teamId}/members", api.GetTeamMembers(teamMemberService)).Methods("GET")
    apiRouter.HandleFunc("/team/{teamId}/member", api.AddTeamMember(teamMemberService)).Methods("POST")
    apiRouter.HandleFunc("/team/{teamId}/member/{userId}", api.RemoveTeamMember(teamMemberService)).Methods("DELETE")

     // Routine routes
     apiRouter.HandleFunc("/routine", api.CreateOrUpdateRoutines(routineService)).Methods("POST")
     apiRouter.HandleFunc("/routine/task/{taskId}", api.GetRoutinesByTaskID(routineService)).Methods("GET")
     apiRouter.HandleFunc("/routine/today/{scheduleType}", api.GetTodayRoutines(routineService)).Methods("GET")
     apiRouter.HandleFunc("/routine/day/{day}/{scheduleType}", api.GetDailyRoutines(routineService)).Methods("GET")
     apiRouter.HandleFunc("/routine/{id}", api.UpdateRoutineDay(routineService)).Methods("PUT")
     apiRouter.HandleFunc("/routine/{id}/status", api.UpdateRoutineStatus(routineService)).Methods("PUT")
     apiRouter.HandleFunc("/routine/task/{taskId}/delete", api.DeleteRoutinesByTaskID(routineService)).Methods("DELETE")
}