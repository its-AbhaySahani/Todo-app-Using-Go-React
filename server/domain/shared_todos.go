package domain

import (
	"time"
)

type SharedTodo struct {
	ID          string
	Task        string
	Description string
	Done        bool
	Important   bool
	UserID      string
	Date        time.Time
	Time        time.Time
	SharedBy    string
}
