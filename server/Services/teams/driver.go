package teams

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/teams_repository"
)

func NewTeamService(repo *teams_repository.TeamRepository) *TeamService {
    return &TeamService{repo: repo}
}