package main

import (
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
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

    // Insert 1000 dummy users
    for i := 0; i < 1000; i++ {
        userID := uuid.New().String()
        username := fmt.Sprintf("user%d", i+1)
        password := "password"
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            log.Fatal("Error hashing password:", err)
        }

        // Check if the user already exists
        var existingUserID string
        err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingUserID)
        if err == nil {
            // User already exists, skip insertion
            continue
        } else if err != sql.ErrNoRows {
            log.Fatal("Error checking user existence:", err)
        }

        _, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)", userID, username, string(hashedPassword))
        if err != nil {
            log.Fatal("Error inserting user:", err)
        }
    }

    fmt.Println("Inserted 1000 dummy users")

    // Fetch all user IDs
    rows, err := db.Query("SELECT id FROM users")
    if err != nil {
        log.Fatal("Error fetching user IDs:", err)
    }
    defer rows.Close()

    var userIDs []string
    for rows.Next() {
        var userID string
        if err := rows.Scan(&userID); err != nil {
            log.Fatal("Error scanning user ID:", err)
        }
        userIDs = append(userIDs, userID)
    }

    // Insert 1000 dummy tasks
    for i := 0; i < 1000; i++ {
        taskID := uuid.New().String()
        task := fmt.Sprintf("Task %d", i+1)
        description := fmt.Sprintf("Description for task %d", i+1)
        done := rand.Intn(2) == 1
        important := rand.Intn(2) == 1
        userID := userIDs[rand.Intn(len(userIDs))]
        date := time.Now().Format("2006-01-02")
        time := time.Now().Format("15:04:05")

        _, err = db.Exec("INSERT INTO todos (id, task, description, done, important, user_id, date, time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", taskID, task, description, done, important, userID, date, time)
        if err != nil {
            log.Fatal("Error inserting task:", err)
        }
    }

    fmt.Println("Inserted 1000 dummy tasks")

    // Insert dummy shared todos
    for i := 0; i < 100; i++ {
        sharedTodoID := uuid.New().String()
        task := fmt.Sprintf("Shared Task %d", i+1)
        description := fmt.Sprintf("Description for shared task %d", i+1)
        done := rand.Intn(2) == 1
        important := rand.Intn(2) == 1
        userID := userIDs[rand.Intn(len(userIDs))]
        sharedBy := userIDs[rand.Intn(len(userIDs))]
        date := time.Now().Format("2006-01-02")
        time := time.Now().Format("15:04:05")

        _, err = db.Exec("INSERT INTO shared_todos (id, task, description, done, important, user_id, date, time, shared_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", sharedTodoID, task, description, done, important, userID, date, time, sharedBy)
        if err != nil {
            log.Fatal("Error inserting shared todo:", err)
        }
    }

    fmt.Println("Inserted 100 dummy shared todos")

    // Insert dummy teams
    for i := 0; i < 10; i++ {
        teamID := uuid.New().String()
        teamName := fmt.Sprintf("Team %d", i+1)
        password := "password"
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            log.Fatal("Error hashing password:", err)
        }
        adminID := userIDs[rand.Intn(len(userIDs))]

        _, err = db.Exec("INSERT INTO teams (id, name, password, admin_id) VALUES (?, ?, ?, ?)", teamID, teamName, string(hashedPassword), adminID)
        if err != nil {
            log.Fatal("Error inserting team:", err)
        }

        // Insert team members
        for j := 0; j < 10; j++ {
            userID := userIDs[rand.Intn(len(userIDs))]
            isAdmin := rand.Intn(2) == 1

            _, err = db.Exec("INSERT INTO team_members (team_id, user_id, is_admin) VALUES (?, ?, ?)", teamID, userID, isAdmin)
            if err != nil {
                log.Fatal("Error inserting team member:", err)
            }
        }

        // Insert dummy team todos
        for k := 0; k < 10; k++ {
            teamTodoID := uuid.New().String()
            task := fmt.Sprintf("Team Task %d", k+1)
            description := fmt.Sprintf("Description for team task %d", k+1)
            done := rand.Intn(2) == 1
            important := rand.Intn(2) == 1
            assignedTo := userIDs[rand.Intn(len(userIDs))]
            date := time.Now().Format("2006-01-02")
            time := time.Now().Format("15:04:05")

            _, err = db.Exec("INSERT INTO team_todos (id, task, description, done, important, team_id, assigned_to, date, time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", teamTodoID, task, description, done, important, teamID, assignedTo, date, time)
            if err != nil {
                log.Fatal("Error inserting team todo:", err)
            }
        }
    }

    fmt.Println("Inserted 10 dummy teams with members and 100 dummy team todos")
}