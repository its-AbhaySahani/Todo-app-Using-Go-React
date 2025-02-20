package domain

import "context"

type TeamMember struct {
    TeamID  string
    UserID  string
    IsAdmin bool
}

type TeamMemberRepository interface {
    AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) error
    GetTeamMembers(ctx context.Context, teamID string) ([]TeamMember, error)
    RemoveTeamMember(ctx context.Context, teamID, userID string) error
}