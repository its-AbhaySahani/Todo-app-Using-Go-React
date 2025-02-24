package todos

import (
    "context"
    "fmt"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/todos_repository"
)

type TodoService struct {
    repo *todos_repository.TodoRepository
}


func (s *TodoService) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.todos.TodoService.CreateTodo"
    res, err := s.repo.CreateTodo(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create todo: %w", functionName, err)
    }
    return res, nil
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    const functionName = "services.todos.TodoService.GetTodosByUserID"
    res, err := s.repo.GetTodosByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get todos by user ID: %w", functionName, err)
    }
    return res, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UpdateTodo"
    res, err := s.repo.UpdateTodo(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update todo: %w", functionName, err)
    }
    return res, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.DeleteTodo"
    res, err := s.repo.DeleteTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete todo: %w", functionName, err)
    }
    return res, nil
}

func (s *TodoService) UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UndoTodo"
    res, err := s.repo.UndoTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to undo todo: %w", functionName, err)
    }
    return res, nil
}