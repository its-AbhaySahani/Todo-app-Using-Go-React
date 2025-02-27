package team_members

import (
    "context"
    "fmt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TeamMemberService struct {
    repo domain.TeamMemberRepository
}

func (s *TeamMemberService) AddTeamMember(ctx context.Context, req *dto.AddTeamMemberRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.team_members.TeamMemberService.AddTeamMember"
    success, err := s.repo.AddTeamMember(ctx, req.TeamID, req.UserID, req.IsAdmin)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to add team member: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TeamMemberService) GetTeamMembers(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
    const functionName = "services.team_members.TeamMemberService.GetTeamMembers"
    domainMembers, err := s.repo.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team members: %w", functionName, err)
    }
    
    // Convert domain.TeamMember to dto.TeamMemberResponse
    var memberResponses []dto.TeamMemberResponse
    for _, member := range domainMembers {
        memberResponses = append(memberResponses, dto.TeamMemberResponse{
            TeamID:  member.TeamID,
            UserID:  member.UserID,
            IsAdmin: member.IsAdmin,
        })
    }
    
    return &dto.TeamMembersResponse{Members: memberResponses}, nil
}

func (s *TeamMemberService) RemoveTeamMember(ctx context.Context, teamID, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.team_members.TeamMemberService.RemoveTeamMember"
    success, err := s.repo.RemoveTeamMember(ctx, teamID, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to remove team member: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}