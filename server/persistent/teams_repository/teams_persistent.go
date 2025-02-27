package teams_repository

import (
    "context"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

// Ensure TeamRepository implements domain.TeamRepository
var _ domain.TeamRepository = (*TeamRepository)(nil)

type TeamRepository struct {
    querier *db.Queries
}

// Implement domain.TeamRepository interface methods
func (r *TeamRepository) CreateTeam(ctx context.Context, name, password, adminID string) (string, error) {
    // Use your existing DTO and converter
    req := &dto.CreateTeamRequest{
        Name:     name,
        Password: password,
        AdminID:  adminID,
    }
    
    params := req.ConvertCreateTeamDomainRequestToPersistentRequest()
    err := r.querier.CreateTeam(ctx, *params)
    if err != nil {
        return "", err
    }
    return params.ID, nil
}

func (r *TeamRepository) GetTeamsByAdminID(ctx context.Context, adminID string) ([]domain.Team, error) {
    teams, err := r.querier.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    
    // Convert db.Team to domain.Team
    domainTeams := make([]domain.Team, len(teams))
    for i, team := range teams {
        domainTeams[i] = domain.Team{
            ID:       team.ID,
            Name:     team.Name,
            Password: team.Password,
            AdminID:  team.AdminID,
        }
    }
    
    return domainTeams, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id string) (domain.Team, error) {
    team, err := r.querier.GetTeamByID(ctx, id)
    if err != nil {
        return domain.Team{}, err
    }
    
    return domain.Team{
        ID:       team.ID,
        Name:     team.Name,
        Password: team.Password,
        AdminID:  team.AdminID,
    }, nil
}

// Original methods for backward compatibility
func (r *TeamRepository) CreateTeamWithDTO(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTeamDomainRequestToPersistentRequest()
    err := r.querier.CreateTeam(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TeamRepository) GetTeamsByAdminIDWithDTO(ctx context.Context, adminID string) (*dto.TeamsResponse, error) {
    teams, err := r.querier.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamsResponse(teams), nil
}

func (r *TeamRepository) GetTeamByIDWithDTO(ctx context.Context, id string) (*dto.TeamResponse, error) {
    team, err := r.querier.GetTeamByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamResponse(&team), nil
}