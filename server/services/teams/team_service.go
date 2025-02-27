package teams

import (
    "context"
    "fmt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TeamService struct {
    repo domain.TeamRepository
}

func (s *TeamService) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateResponse, error) {
    const functionName = "services.teams.TeamService.CreateTeam"
    id, err := s.repo.CreateTeam(ctx, req.Name, req.Password, req.AdminID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create team: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *TeamService) GetTeamsByAdminID(ctx context.Context, adminID string) (*dto.TeamsResponse, error) {
    const functionName = "services.teams.TeamService.GetTeamsByAdminID"
    domainTeams, err := s.repo.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get teams by admin ID: %w", functionName, err)
    }
    
    // Convert domain.Team to dto.TeamResponse
    var teamResponses []dto.TeamResponse
    for _, team := range domainTeams {
        teamResponses = append(teamResponses, dto.TeamResponse{
            ID:       team.ID,
            Name:     team.Name,
            Password: team.Password,
            AdminID:  team.AdminID,
        })
    }
    
    return &dto.TeamsResponse{Teams: teamResponses}, nil
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*dto.TeamResponse, error) {
    const functionName = "services.teams.TeamService.GetTeamByID"
    team, err := s.repo.GetTeamByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get team by ID: %w", functionName, err)
    }
    
    return &dto.TeamResponse{
        ID:       team.ID,
        Name:     team.Name,
        Password: team.Password,
        AdminID:  team.AdminID,
    }, nil
}