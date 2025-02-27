package team_members_repository

import (
    "context"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure TeamMemberRepository implements domain.TeamMemberRepository
var _ domain.TeamMemberRepository = (*TeamMemberRepository)(nil)

type TeamMemberRepository struct {
    querier *db.Queries
}

// Implement domain.TeamMemberRepository interface methods
func (r *TeamMemberRepository) AddTeamMember(ctx context.Context, teamID, userID string, isAdmin bool) (bool, error) {
    // Use your existing DTO and converter
    req := &dto.AddTeamMemberRequest{
        TeamID:  teamID,
        UserID:  userID,
        IsAdmin: isAdmin,
    }
    
    params := req.ConvertAddTeamMemberDomainRequestToPersistentRequest()
    err := r.querier.AddTeamMember(ctx, *params)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *TeamMemberRepository) GetTeamMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error) {
    members, err := r.querier.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, err
    }
    
    // Convert db.TeamMember to domain.TeamMember
    domainMembers := make([]domain.TeamMember, len(members))
    for i, member := range members {
        domainMembers[i] = domain.TeamMember{
            TeamID:  member.TeamID,
            UserID:  member.UserID,
            IsAdmin: member.IsAdmin.Bool,
        }
    }
    
    return domainMembers, nil
}

func (r *TeamMemberRepository) RemoveTeamMember(ctx context.Context, teamID, userID string) (bool, error) {
    err := r.querier.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
        TeamID: teamID,
        UserID: userID,
    })
    if err != nil {
        return false, err
    }
    return true, nil
}

// Original methods for backward compatibility
func (r *TeamMemberRepository) AddTeamMemberWithDTO(ctx context.Context, req *dto.AddTeamMemberRequest) (*dto.SuccessResponse, error) {
    params := req.ConvertAddTeamMemberDomainRequestToPersistentRequest()
    err := r.querier.AddTeamMember(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}

func (r *TeamMemberRepository) GetTeamMembersWithDTO(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
    members, err := r.querier.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamMembersResponse(members), nil
}

func (r *TeamMemberRepository) RemoveTeamMemberWithDTO(ctx context.Context, teamID, userID string) (*dto.SuccessResponse, error) {
    err := r.querier.RemoveTeamMember(ctx, db.RemoveTeamMemberParams{
        TeamID: teamID,
        UserID: userID,
    })
    if err != nil {
        return nil, err
    }
    return &dto.SuccessResponse{Success: true}, nil
}