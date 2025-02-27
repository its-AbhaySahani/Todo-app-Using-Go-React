package teams

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamService(repo domain.TeamRepository) *TeamService {
    return &TeamService{repo: repo}
}