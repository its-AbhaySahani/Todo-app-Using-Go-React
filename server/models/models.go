package models

import (
    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/Database"
)

type Todo struct {
    ID   string `json:"id"`
    Task string `json:"task"`
    Done bool   `json:"done"`
}

func GetTodos() ([]Todo, error) {
    rows, err := database.DB.Query("SELECT id, task, done FROM todos")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }
    return todos, nil
}

func CreateTodo(task string) (Todo, error) {
    id := uuid.New().String()
    _, err := database.DB.Exec("INSERT INTO todos (id, task, done) VALUES (?, ?, ?)", id, task, false)
    if err != nil {
        return Todo{}, err
    }
    return Todo{ID: id, Task: task, Done: false}, nil
}

func UpdateTodo(id string, task string, done bool) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET task = ?, done = ? WHERE id = ?", task, done, id)
    if err != nil {
        return Todo{}, err
    }
    return Todo{ID: id, Task: task, Done: done}, nil
}

func DeleteTodo(id string) error {
    _, err := database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
    return err
}

func UndoTodo(id string) (Todo, error) {
    _, err := database.DB.Exec("UPDATE todos SET done = ? WHERE id = ?", false, id)
    if err != nil {
        return Todo{}, err
    }
    return Todo{ID: id, Done: false}, nil
}