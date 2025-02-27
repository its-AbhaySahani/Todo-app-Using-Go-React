package shared_todos

import (
    "context"
    "fmt"
    
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/dto"
)

type SharedTodoService struct {
    repo domain.SharedTodoRepository
}

func (s *SharedTodoService) CreateSharedTodo(ctx context.Context, req *dto.CreateSharedTodoRequest) (*dto.CreateResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.CreateSharedTodo"
    id, err := s.repo.CreateSharedTodo(ctx, req.Task, req.Description, req.Done, req.Important, req.UserID, req.SharedBy)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to create shared todo: %w", functionName, err)
    }
    return &dto.CreateResponse{ID: id}, nil
}

func (s *SharedTodoService) GetSharedTodos(ctx context.Context, userID string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedTodos"
    domainTodos, err := s.repo.GetSharedTodos(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared todos: %w", functionName, err)
    }
    
    // Convert domain.SharedTodo to dto.SharedTodoResponse
    var receivedTodos []dto.SharedTodoResponse
    for _, todo := range domainTodos {
        receivedTodos = append(receivedTodos, dto.SharedTodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
            SharedBy:    todo.SharedBy,
        })
    }
    
    return &dto.SharedTodosResponse{Received: receivedTodos}, nil
}

func (s *SharedTodoService) GetSharedByMeTodos(ctx context.Context, sharedBy string) (*dto.SharedTodosResponse, error) {
    const functionName = "services.shared_todos.SharedTodoService.GetSharedByMeTodos"
    domainTodos, err := s.repo.GetSharedByMeTodos(ctx, sharedBy)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to get shared by me todos: %w", functionName, err)
    }
    
    // Convert domain.SharedTodo to dto.SharedTodoResponse
    var sharedTodos []dto.SharedTodoResponse
    for _, todo := range domainTodos {
        sharedTodos = append(sharedTodos, dto.SharedTodoResponse{
            ID:          todo.ID,
            Task:        todo.Task,
            Description: todo.Description,
            Done:        todo.Done,
            Important:   todo.Important,
            UserID:      todo.UserID,
            Date:        todo.Date,
            Time:        todo.Time,
            SharedBy:    todo.SharedBy,
        })
    }
    
    return &dto.SharedTodosResponse{Shared: sharedTodos}, nil
}