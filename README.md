# Todo-app-Using-Go-React

<!-- package database

import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
    dsn := "Abhay:Abhay@123@tcp(mysql:3306)/Todo_app"
    var db *sql.DB
    var err error

    for i := 0; i < 10; i++ {
        db, err = sql.Open("mysql", dsn)
        if err == nil {
            err = db.Ping()
            if err == nil {
                break
            }
        }
        fmt.Println("Error connecting to the database. Retrying in 5 seconds...")
        time.Sleep(5 * time.Second)
    }

    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        return
    }

    DB = db
    fmt.Println("Connected to the database successfully")
} -->


## To run the updated docker app
docker-compose up --build