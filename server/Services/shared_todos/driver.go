package shared_todos

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/shared_todos_repository"
)

func NewSharedTodoService(repo *shared_todos_repository.SharedTodoRepository) *SharedTodoService {
    return &SharedTodoService{repo: repo}
}