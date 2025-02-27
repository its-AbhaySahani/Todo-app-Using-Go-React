package domain

import (
    "context"
)

type Team struct {
    ID       string
    Name     string
    Password string
    AdminID  string
}

// TeamRepository defines the interface for team persistence operations
type TeamRepository interface {
    CreateTeam(ctx context.Context, name, password, adminID string) (string, error)
    GetTeamsByAdminID(ctx context.Context, adminID string) ([]Team, error)
    GetTeamByID(ctx context.Context, id string) (Team, error)
}