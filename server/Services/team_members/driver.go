package team_members

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/team_members_repository"
)

func NewTeamMemberService(repo *team_members_repository.TeamMemberRepository) *TeamMemberService {
    return &TeamMemberService{repo: repo}
}