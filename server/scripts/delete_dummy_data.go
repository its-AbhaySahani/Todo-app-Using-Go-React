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

    // Delete dummy users
    _, err = db.Exec("DELETE FROM users WHERE username LIKE 'user%'")
    if err != nil {
        log.Fatal("Error deleting users:", err)
    }

    fmt.Println("Deleted dummy users")
}