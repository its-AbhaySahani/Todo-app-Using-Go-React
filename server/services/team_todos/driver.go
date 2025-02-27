package team_todos

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamTodoService(repo domain.TeamTodoRepository) *TeamTodoService {
    return &TeamTodoService{repo: repo}
}