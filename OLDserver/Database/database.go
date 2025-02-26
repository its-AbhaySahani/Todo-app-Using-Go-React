package database

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
    dsn := "Abhay:Abhay@123@tcp(127.0.0.1:3306)/Todo_app"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal("Error pinging the database:", err)
    }

    DB = db
    fmt.Println("Connected to the database successfully")
}