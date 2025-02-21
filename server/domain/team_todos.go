package domain

import (
    "time"
)

type TeamTodo struct {
    ID          string
    Task        string
    Description string
    Done        bool
    Important   bool
    TeamID      string
    AssignedTo  string
    Date        time.Time
    Time        time.Time
}
