# Todo-app-Using-Go-React
<!-- -- Users Queries

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

-- name: UpdateTodo :exec
UPDATE todos
SET task = ?, description = ?, done = ?, important = ?
WHERE id = ? AND user_id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ? AND user_id = ?;

-- name: UndoTodo :exec
UPDATE todos
SET done = false
WHERE id = ? AND user_id = ?;

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

-- name: ShareTodoWithUser :exec
INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by)
SELECT ?, task, description, done, important, (SELECT id FROM users WHERE username = ?), date, time, ?
FROM todos
WHERE todos.id = ?;

-- Teams Queries

-- name: CreateTeam :exec
INSERT INTO teams (id, name, password, admin_id)
VALUES (?, ?, ?, ?);

-- name: GetTeamsByAdminID :many
SELECT id, name, password, admin_id
FROM teams
WHERE admin_id = ?;

-- name: GetTeamByID :one
SELECT id, name, password, admin_id
FROM teams
WHERE id = ?;

-- Team Members Queries

-- name: AddTeamMember :exec
INSERT INTO team_members (team_id, user_id, is_admin)
VALUES (?, ?, ?);

-- name: GetTeamMembers :many
SELECT team_id, user_id, is_admin
FROM team_members
WHERE team_id = ?;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE team_id = ? AND user_id = ?;

-- Team Todos Queries

-- name: CreateTeamTodo :exec
INSERT INTO team_todos (id, task, description, done, important, team_id, assigned_to, date, time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTeamTodos :many
SELECT id, task, description, done, important, team_id, assigned_to, date, time
FROM team_todos
WHERE team_id = ?;

-- name: UpdateTeamTodo :exec
UPDATE team_todos
SET task = ?, description = ?, done = ?, important = ?, assigned_to = ?
WHERE id = ? AND team_id = ?;

-- name: DeleteTeamTodo :exec
DELETE FROM team_todos
WHERE id = ? AND team_id = ?;

-- name: JoinTeam :exec
INSERT INTO team_members (team_id, user_id, is_admin)
SELECT id, ?, false
FROM teams
WHERE name = ? AND password = ?; -->


## To run the updated docker app
docker-compose up --build


