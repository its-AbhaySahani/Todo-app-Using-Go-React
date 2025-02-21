package teams_repository

import (
    "context"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

type TeamServiceRepository interface {
    CreateTeam(ctx context.Context, name, password, adminID string) error
    GetTeamsByAdminID(ctx context.Context, adminID string) ([]domain.Team, error)
    GetTeamByID(ctx context.Context, id string) (domain.Team, error)
}
type TeamRepository struct {
    querier *db.Queries
}

func (r *TeamRepository) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateResponse, error) {
    params := req.ConvertCreateTeamDomainRequestToPersistentRequest()
    err := r.querier.CreateTeam(ctx, *params)
    if err != nil {
        return nil, err
    }
    return &dto.CreateResponse{ID: params.ID}, nil
}

func (r *TeamRepository) GetTeamsByAdminID(ctx context.Context, adminID string) (*dto.TeamsResponse, error) {
    teams, err := r.querier.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamsResponse(teams), nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id string) (*dto.TeamResponse, error) {
    team, err := r.querier.GetTeamByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return dto.NewTeamResponse(&team), nil
}

