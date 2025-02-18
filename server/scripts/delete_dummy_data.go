package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "Abhay:Abhay@123@tcp(127.0.0.1:3306)/Todo_app"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Error pinging the database:", err)
    }

    fmt.Println("Connected to the database successfully")

    // Delete dummy tasks
    _, err = db.Exec("DELETE FROM todos WHERE task LIKE 'Task %'")
    if err != nil {
        log.Fatal("Error deleting tasks:", err)
    }

    fmt.Println("Deleted dummy tasks")

    // Delete dummy shared todos
    _, err = db.Exec("DELETE FROM shared_todos WHERE task LIKE 'Shared Task %'")
    if err != nil {
        log.Fatal("Error deleting shared todos:", err)
    }

    fmt.Println("Deleted dummy shared todos")

    // Delete dummy team todos
    _, err = db.Exec("DELETE FROM team_todos WHERE task LIKE 'Team Task %'")
    if err != nil {
        log.Fatal("Error deleting team todos:", err)
    }

    fmt.Println("Deleted dummy team todos")

    // Delete dummy team members
    _, err = db.Exec("DELETE FROM team_members WHERE team_id IN (SELECT id FROM teams WHERE name LIKE 'Team %')")
    if err != nil {
        log.Fatal("Error deleting team members:", err)
    }

    fmt.Println("Deleted dummy team members")

    // Delete dummy teams
    _, err = db.Exec("DELETE FROM teams WHERE name LIKE 'Team %'")
    if err != nil {
        log.Fatal("Error deleting teams:", err)
    }

    fmt.Println("Deleted dummy teams")

    // Delete dummy users
    _, err = db.Exec("DELETE FROM users WHERE username LIKE 'user%'")
    if err != nil {
        log.Fatal("Error deleting users:", err)
    }

    fmt.Println("Deleted dummy users")
}