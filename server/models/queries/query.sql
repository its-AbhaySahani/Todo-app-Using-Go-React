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

-- name: GetTeams :many
SELECT t.id, t.name, t.password, t.admin_id
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = ? /* sqlc.arg(userID) */;

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

-- name: GetTeamMemberDetails :many
SELECT u.id, u.username, tm.is_admin
FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = ? /* sqlc.arg(teamID) */;

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

-- Routines Queries (new additions for the routines functionality)

-- name: CreateRoutine :exec
INSERT INTO routines (id, day, scheduleType, taskId, userId, createdAt, updatedAt, isActive)
VALUES (
  ? /* sqlc.arg(id) */,
  ? /* sqlc.arg(day) */,
  ? /* sqlc.arg(scheduleType) */,
  ? /* sqlc.arg(taskId) */,
  ? /* sqlc.arg(userId) */,
  ? /* sqlc.arg(createdAt) */,
  ? /* sqlc.arg(updatedAt) */,
  ? /* sqlc.arg(isActive) */
);

-- name: UpdateRoutineStatus :exec
UPDATE routines
SET isActive = ? /* sqlc.arg(isActive) */,
    updatedAt = ? /* sqlc.arg(updatedAt) */
WHERE id = ? /* sqlc.arg(id) */;

-- name: UpdateRoutineDay :exec
UPDATE routines
SET day = ? /* sqlc.arg(day) */,
    updatedAt = ? /* sqlc.arg(updatedAt) */
WHERE id = ? /* sqlc.arg(id) */;

-- name: GetRoutinesByTaskID :many
SELECT id, day, scheduleType, taskId, userId, createdAt, updatedAt, isActive
FROM routines
WHERE taskId = ? /* sqlc.arg(taskId) */;

-- name: GetDailyRoutines :many
SELECT t.id, t.task, t.description, t.done, t.important, t.user_id, t.date, t.time
FROM todos t
JOIN routines r ON t.id = r.taskId
WHERE r.day = ? /* sqlc.arg(day) */ 
  AND r.scheduleType = ? /* sqlc.arg(scheduleType) */ 
  AND r.userId = ? /* sqlc.arg(userId) */ 
  AND r.isActive = true;

-- name: DeleteRoutinesByTaskID :exec
DELETE FROM routines
WHERE taskId = ? /* sqlc.arg(taskId) */;