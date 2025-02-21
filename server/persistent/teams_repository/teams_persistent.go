package teams_repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/models/db"
)

type TeamRepository struct {
    querier *db.Queries
}

func (r *TeamRepository) CreateTeam(ctx context.Context, name, password, adminID string) error {
    id := uuid.New().String()
    return r.querier.CreateTeam(ctx, db.CreateTeamParams{
        ID:       id,
        Name:     name,
        Password: password,
        AdminID:  adminID,
    })
}

func (r *TeamRepository) GetTeamsByAdminID(ctx context.Context, adminID string) ([]domain.Team, error) {
    teams, err := r.querier.GetTeamsByAdminID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    var result []domain.Team
    for _, team := range teams {
        result = append(result, domain.Team{
            ID:       team.ID,
            Name:     team.Name,
            Password: team.Password,
            AdminID:  team.AdminID,
        })
    }
    return result, nil
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