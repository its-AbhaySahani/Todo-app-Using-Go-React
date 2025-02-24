package teams

import (
    "context"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/teams_repository"
)

type TeamService struct {
    repo *teams_repository.TeamRepository
}


func (s *TeamService) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateResponse, error) {
    const functionName = "services.teams.TeamService.CreateTeam"
    res, err := s.repo.CreateTeam(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create team: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamService) GetTeamsByAdminID(ctx context.Context, adminID string) (*dto.TeamsResponse, error) {
    const functionName = "services.teams.TeamService.GetTeamsByAdminID"
    res, err := s.repo.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get teams by admin ID: %w", functionName, err)
    }
    return res, nil
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*dto.TeamResponse, error) {
    const functionName = "services.teams.TeamService.GetTeamByID"
    res, err := s.repo.GetTeamByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team by ID: %w", functionName, err)
    }
    return res, nil
}

// func (s *TeamService) JoinTeam(ctx context.Context, teamName, password, userID string) (*dto.SuccessResponse, error) {
//     const functionName = "services.teams.TeamService.JoinTeam"
//     res, err := s.repo.JoinTeam(ctx, teamName, password, userID)
//     if err != nil {
//         return nil, fmt.Errorf("%s: failed to join team: %w", functionName, err)
//     }
//     return res, nil
// }

// func (s *TeamService) GetTeamDetails(ctx context.Context, teamID string) (*dto.TeamResponse, error) {
//     const functionName = "services.teams.TeamService.GetTeamDetails"
//     res, err := s.repo.GetTeamDetails(ctx, teamID)
//     if err != nil {
//         return nil, fmt.Errorf("%s: failed to get team details: %w", functionName, err)
//     }
//     return res, nil
// }

// func (s *TeamService) GetTeamMembers(ctx context.Context, teamID string) (*dto.TeamMembersResponse, error) {
//     const functionName = "services.teams.TeamService.GetTeamMembers"
//     res, err := s.repo.GetTeamMembers(ctx, teamID)
//     if err != nil {
//         return nil, fmt.Errorf("%s: failed to get team members: %w", functionName, err)
//     }
//     return res, nil
// }

// func (s *TeamService) RemoveTeamMember(ctx context.Context, teamID, userID string) (*dto.SuccessResponse, error) {
//     const functionName = "services.teams.TeamService.RemoveTeamMember"
//     res, err := s.repo.RemoveTeamMember(ctx, teamID, userID)
//     if err != nil {
//         return nil, fmt.Errorf("%s: failed to remove team member: %w", functionName, err)
//     }
//     return res, nil
// }