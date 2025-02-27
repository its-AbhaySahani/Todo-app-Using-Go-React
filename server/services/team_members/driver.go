package team_members

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTeamMemberService(repo domain.TeamMemberRepository) *TeamMemberService {
    return &TeamMemberService{repo: repo}
}