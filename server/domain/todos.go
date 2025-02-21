package domain

import (
    "time"
)

type Todo struct {
    ID          string
    Task        string
    Description string
    Done        bool
    Important   bool
    UserID      string
    Date        time.Time
    Time        time.Time
}
