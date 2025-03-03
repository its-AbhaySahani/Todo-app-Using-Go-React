package routines

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewRoutineService(repo domain.RoutineRepository) *RoutineService {
    return &RoutineService{repo: repo}
}