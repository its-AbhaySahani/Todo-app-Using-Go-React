package shared_todos

import (
	"github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewSharedTodoService(repo domain.SharedTodoRepository, todoRepo domain.TodoRepository, userRepo domain.UserRepository) *SharedTodoService {
    return &SharedTodoService{
        repo:     repo,
        todoRepo: todoRepo,
        userRepo: userRepo,
    }
}