package todos

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/todos_repository"
)

func NewTodoService(repo *todos_repository.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}