package todos

import (
    "context"
    "fmt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type TodoService struct {
    repo domain.TodoRepository
}

func (s *TodoService) CreateTodo(ctx context.Context, req *dto.CreateTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.todos.TodoService.CreateTodo"
    id, err := s.repo.CreateTodo(ctx, req.Task, req.Description, req.Done, req.Important, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create todo: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userID string) (*dto.TodosResponse, error) {
    const functionName = "services.todos.TodoService.GetTodosByUserID"
    domainTodos, err := s.repo.GetTodosByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get todos by user ID: %w", functionName, err)
    }
    
    // Convert domain.Todo to dto.TodoResponse
    var todoResponses []dto.TodoResponse
    for _, todo := range domainTodos {
        todoResponses = append(todoResponses, dto.TodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
        })
    }
    
    return &dto.TodosResponse{Todos: todoResponses}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, req *dto.UpdateTodoRequest) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UpdateTodo"
    success, err := s.repo.UpdateTodo(ctx, req.ID, req.Task, req.Description, req.Done, req.Important, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to update todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.DeleteTodo"
    success, err := s.repo.DeleteTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to delete todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}

func (s *TodoService) UndoTodo(ctx context.Context, id, userID string) (*dto.SuccessResponse, error) {
    const functionName = "services.todos.TodoService.UndoTodo"
    success, err := s.repo.UndoTodo(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to undo todo: %w", functionName, err)
    }
    return &dto.SuccessResponse{Success: success}, nil
}