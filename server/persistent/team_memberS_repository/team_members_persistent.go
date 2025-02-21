package team_members_repository

import (
    "context"
    "database/sql"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type TeamMemberRepository struct {
    querier *db.Queries
}

func (r *TeamMemberRepository) AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) error {
    return r.querier.AddTeamMember(ctx, db.AddTeamMemberParams{
        TeamID:  teamID,
        UserID:  userID,
        IsAdmin: sql.NullBool{Bool: isAdmin, Valid: true},
    })
}

func (r *TeamMemberRepository) GetTeamMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error) {
    members, err := r.querier.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, err
    }
    var result []domain.TeamMember
    for _, member := range members {
        result = append(result, domain.TeamMember{
            TeamID:  member.TeamID,
            UserID:  member.UserID,
            IsAdmin: member.IsAdmin.Bool,
        })
    }
    return result, nil
}

func (r *TeamMemberRepository) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
    return r.querier.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
        TeamID: teamID,
        UserID: userID,
    })
}