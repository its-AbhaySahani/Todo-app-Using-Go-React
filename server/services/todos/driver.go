package todos

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewTodoService(repo domain.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}