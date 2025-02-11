package models

import (
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
    "time"
)

type Todo struct {
    ID    string `json:"id"`
    Task  string `json:"task"`
    Done  bool   `json:"done"`
    UserID string `json:"user_id"`
    Date  string `json:"date"`
    Time  string `json:"time"`
}

func GetTodos(userID string) ([]Todo, error) {
    rows, err := database.DB.Query("SELECT id, task, done, date, time FROM todos WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done, &todo.Date, &todo.Time); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }
    return todos, nil
}

func CreateTodo(task string, userID string) (Todo, error) {
    id := uuid.New().String()
    currentDate := time.Now().Format("2006-01-02")
    currentTime := time.Now().Format("15:04:05")
    _, err := database.DB.Exec("INSERT INTO todos (id, task, done, user_id, date, time) VALUES (?, ?, ?, ?, ?, ?)", id, task, false, userID, currentDate, currentTime)
    if err != nil {
        return Todo{}, err
    }
    return Todo{ID: id, Task: task, Done: false, UserID: userID, Date: currentDate, Time: currentTime}, nil
}

func UpdateTodo(id string, task string, done bool, userID string) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET task = ?, done = ? WHERE id = ? AND user_id = ?", task, done, id, userID)
    if err != nil {
        return Todo{}, err
    }
    var todo Todo
    err = database.DB.QueryRow("SELECT id, task, done, date, time FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&todo.ID, &todo.Task, &todo.Done, &todo.Date, &todo.Time)
    if err != nil {
        return Todo{}, err
    }
    return todo, nil
}

func DeleteTodo(id string, userID string) error {
    _, err := database.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", id, userID)
    return err
}

func UndoTodo(id string, userID string) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET done = ? WHERE id = ? AND user_id = ?", false, id, userID)
    if err != nil {
        return Todo{}, err
    }
    var todo Todo
    err = database.DB.QueryRow("SELECT id, task, done, date, time FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&todo.ID, &todo.Task, &todo.Done, &todo.Date, &todo.Time)
    if err != nil {
        return Todo{}, err
    }
    return todo, nil
}