package shared_todos

import (
    "context"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/shared_todos_repository"
)

type SharedTodoService struct {
    repo *shared_todos_repository.SharedTodoRepository
}

func (s *SharedTodoService) CreateSharedTodo(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.CreateSharedTodo"
    res, err := s.repo.CreateSharedTodo(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create shared todo: %w", functionName, err)
    }
    return res, nil
}

func (s *SharedTodoService) GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedTodos"
    res, err := s.repo.GetSharedTodos(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared todos: %w", functionName, err)
    }
    return res, nil
}

func (s *SharedTodoService) GetSharedByMeTodos(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedByMeTodos"
    res, err := s.repo.GetSharedByMeTodos(ctx, sharedBy)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared by me todos: %w", functionName, err)
    }
    return res, nil
}