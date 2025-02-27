package shared_todos

import (
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewSharedTodoService(repo domain.SharedTodoRepository) *SharedTodoService {
    return &SharedTodoService{repo: repo}
}