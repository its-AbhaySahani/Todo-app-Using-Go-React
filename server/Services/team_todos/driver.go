package team_todos

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_todos_repository"
)

func NewTeamTodoService(repo *team_todos_repository.TeamTodoRepository) *TeamTodoService {
    return &TeamTodoService{repo: repo}
}