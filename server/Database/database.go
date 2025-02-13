package database

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
    dsn := "Abhay:Abhay@123@tcp(127.0.0.1:3306)/Todo_app"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        return
    }

    if err := db.Ping(); err != nil {
        fmt.Println("Error pinging the database:", err)
        return
    }

    DB = db
    fmt.Println("Connected to the database successfully")
}
