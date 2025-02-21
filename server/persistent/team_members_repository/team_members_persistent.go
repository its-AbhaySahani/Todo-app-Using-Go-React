package team_members_repository

import (
    "context"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

type TeamMemberServiceRepository interface {
    AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) error
    GetTeamMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error)
    RemoveTeamMember(ctx context.Context, teamID, userID string) error
}
type TeamMemberRepository struct {
    querier *db.Queries
}

func (r *TeamMemberRepository) AddTeamMember(ctx context.Context, req *dto.AddTeamMemberRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertAddTeamMemberDomainRequestToPersistentRequest()
    err := r.querier.AddTeamMember(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TeamMemberRepository) GetTeamMembers(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
    members, err := r.querier.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamMembersResponse(members), nil
}

func (r *TeamMemberRepository) RemoveTeamMember(ctx context.Context, teamID, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
        TeamID: teamID,
        UserID: userID,
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}