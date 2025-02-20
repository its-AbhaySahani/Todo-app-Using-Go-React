-- Users Queries

-- name: CreateUser :exec
INSERT INTO users (id, username, password)
VALUES (?, ?, ?);

-- name: GetUserByUsername :one
SELECT id, username, password
FROM users
WHERE username = ?;

-- Todos Queries

-- name: CreateTodo :exec
INSERT INTO todos (id, task, description, done, important, user_id, date, time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTodosByUserID :many
SELECT id, task, description, done, important, user_id, date, time
FROM todos
WHERE user_id = ?;

-- Shared Todos Queries

-- name: CreateSharedTodo :exec
INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetSharedTodos :many
SELECT id, task, description, done, important, user_id, date, time, shared_by
FROM shared_todos
WHERE user_id = ?;

-- name: GetSharedByMeTodos :many
SELECT id, task, description, done, important, user_id, date, time, shared_by
FROM shared_todos
WHERE shared_by = ?;

-- Teams Queries

-- name: CreateTeam :exec
INSERT INTO teams (id, name, password, admin_id)
VALUES (?, ?, ?, ?);

-- name: GetTeamsByAdminID :many
SELECT id, name, password, admin_id
FROM teams
WHERE admin_id = ?;

-- Team Members Queries

-- name: AddTeamMember :exec
INSERT INTO team_members (team_id, user_id, is_admin)
VALUES (?, ?, ?);

-- name: GetTeamMembers :many
SELECT team_id, user_id, is_admin
FROM team_members
WHERE team_id = ?;

-- Team Todos Queries

-- name: CreateTeamTodo :exec
INSERT INTO team_todos (id, task, description, done, important, team_id, assigned_to, date, time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTeamTodos :many
SELECT id, task, description, done, important, team_id, assigned_to, date, time
FROM team_todos
WHERE team_id = ?;