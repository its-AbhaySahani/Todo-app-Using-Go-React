package team_members

import (
    "context"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_members_repository"
)

type TeamMemberService struct {
    repo *team_members_repository.TeamMemberRepository
}


func (s *TeamMemberService) AddTeamMember(ctx context.Context, req *dto.AddTeamMemberRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.team_members.TeamMemberService.AddTeamMember"
    res, err := s.repo.AddTeamMember(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to add team member: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamMemberService) GetTeamMembers(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
    const functionName = "services.team_members.TeamMemberService.GetTeamMembers"
    res, err := s.repo.GetTeamMembers(ctx, teamID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team members: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamMemberService) RemoveTeamMember(ctx context.Context, teamID, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.team_members.TeamMemberService.RemoveTeamMember"
    res, err := s.repo.RemoveTeamMember(ctx, teamID, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to remove team member: %w", functionName, err)
    }
    return res, nil
}