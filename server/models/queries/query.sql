-- Users Queries

-- name: CreateUser :exec
INSERT INTO users (id, username, password)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(username) */,
  ? /* sqlc.arg(password) */
);

-- name: GetUserByUsername :one
SELECT id, username, password
FROM users
WHERE username = ? /* sqlc.arg(username) */;

-- Todos Queries

-- name: CreateTodo :exec
INSERT INTO todos (id, task, description, done, important, user_id, date, time)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(task) */,
  ? /* sqlc.arg(description) */,
  ? /* sqlc.arg(done) */,
  ? /* sqlc.arg(important) */,
  ? /* sqlc.arg(userID) */,
  ? /* sqlc.arg(date) */,
  ? /* sqlc.arg(time) */
);

-- name: GetTodosByUserID :many
SELECT id, task, description, done, important, user_id, date, time
FROM todos
WHERE user_id = ? /* sqlc.arg(userID) */;

-- name: UpdateTodo :exec
UPDATE todos
SET task = ? /* sqlc.arg(task) */,
    description = ? /* sqlc.arg(description) */,
    done = ? /* sqlc.arg(done) */,
    important = ? /* sqlc.arg(important) */
WHERE id = ? /* sqlc.arg(id) */ AND user_id = ? /* sqlc.arg(userID) */;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ? /* sqlc.arg(id) */ AND user_id = ? /* sqlc.arg(userID) */;

-- name: UndoTodo :exec
UPDATE todos
SET done = false
WHERE id = ? /* sqlc.arg(id) */ AND user_id = ? /* sqlc.arg(userID) */;

-- Shared Todos Queries

-- name: CreateSharedTodo :exec
INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(task) */,
  ? /* sqlc.arg(description) */,
  ? /* sqlc.arg(done) */,
  ? /* sqlc.arg(important) */,
  ? /* sqlc.arg(userID) */,
  ? /* sqlc.arg(date) */,
  ? /* sqlc.arg(time) */,
  ? /* sqlc.arg(sharedBy) */
);

-- name: GetSharedTodos :many
SELECT id, task, description, done, important, user_id, date, time, shared_by
FROM shared_todos
WHERE user_id = ? /* sqlc.arg(userID) */;

-- name: GetSharedByMeTodos :many
SELECT id, task, description, done, important, user_id, date, time, shared_by
FROM shared_todos
WHERE shared_by = ? /* sqlc.arg(sharedBy) */;

-- name: ShareTodoWithUser :exec
INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by)
SELECT 
  ? /* sqlc.arg(newID) */,
  task, 
  description, 
  done, 
  important, 
  (SELECT id FROM users WHERE username = ? /* sqlc.arg(receiverUsername) */), 
  date, 
  time, 
  ? /* sqlc.arg(senderID) */
FROM todos
WHERE todos.id = ? /* sqlc.arg(todoID) */;

-- Teams Queries

-- name: CreateTeam :exec
INSERT INTO teams (id, name, password, admin_id)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(name) */,
  ? /* sqlc.arg(password) */,
  ? /* sqlc.arg(adminID) */
);

-- name: GetTeamsByAdminID :many
SELECT id, name, password, admin_id
FROM teams
WHERE admin_id = ? /* sqlc.arg(adminID) */;

-- name: GetTeamByID :one
SELECT id, name, password, admin_id
FROM teams
WHERE id = ? /* sqlc.arg(teamID) */;

-- Team Members Queries

-- name: AddTeamMember :exec
INSERT INTO team_members (team_id, user_id, is_admin)
VALUES (
  ? /* sqlc.arg(teamID) */,
  ? /* sqlc.arg(userID) */,
  ? /* sqlc.arg(isAdmin) */
);

-- name: GetTeamMembers :many
SELECT team_id, user_id, is_admin
FROM team_members
WHERE team_id = ? /* sqlc.arg(teamID) */;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE team_id = ? /* sqlc.arg(teamID) */ AND user_id = ? /* sqlc.arg(userID) */;

-- Team Todos Queries

-- name: CreateTeamTodo :exec
INSERT INTO team_todos (id, task, description, done, important, team_id, assigned_to, date, time)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(task) */,
  ? /* sqlc.arg(description) */,
  ? /* sqlc.arg(done) */,
  ? /* sqlc.arg(important) */,
  ? /* sqlc.arg(teamID) */,
  ? /* sqlc.arg(assignedTo) */,
  ? /* sqlc.arg(date) */,
  ? /* sqlc.arg(time) */
);

-- name: GetTeamTodos :many
SELECT id, task, description, done, important, team_id, assigned_to, date, time
FROM team_todos
WHERE team_id = ? /* sqlc.arg(teamID) */;

-- name: UpdateTeamTodo :exec
UPDATE team_todos
SET 
  task = ? /* sqlc.arg(task) */,
  description = ? /* sqlc.arg(description) */,
  done = ? /* sqlc.arg(done) */,
  important = ? /* sqlc.arg(important) */,
  assigned_to = ? /* sqlc.arg(assignedTo) */
WHERE id = ? /* sqlc.arg(id) */ AND team_id = ? /* sqlc.arg(teamID) */;

-- name: DeleteTeamTodo :exec
DELETE FROM team_todos
WHERE id = ? /* sqlc.arg(id) */ AND team_id = ? /* sqlc.arg(teamID) */;

-- name: JoinTeam :exec
INSERT INTO team_members (team_id, user_id, is_admin)
SELECT 
  id, 
  ? /* sqlc.arg(userID) */, 
  false
FROM teams
WHERE name = ? /* sqlc.arg(teamName) */ AND password = ? /* sqlc.arg(teamPassword) */;