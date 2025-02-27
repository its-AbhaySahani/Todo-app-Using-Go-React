package domain

import (
    "context"
)

type TeamMember struct {
    TeamID  string
    UserID  string
    IsAdmin bool
}

// TeamMemberRepository defines the interface for team member persistence operations
type TeamMemberRepository interface {
    AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) (bool, error)
    GetTeamMembers(ctx context.Context, teamID string) ([]TeamMember, error)
    RemoveTeamMember(ctx context.Context, teamID, userID string) (bool, error)
}