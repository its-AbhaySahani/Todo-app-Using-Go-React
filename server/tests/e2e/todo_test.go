package e2e

import (
    "testing"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler/dto"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/tests/e2e/helpers"
    "github.com/stretchr/testify/suite"
)

type TodoE2ETestSuite struct {
    E2ETestSuite
    token string
}

func TestTodoE2E(t *testing.T) {
    suite.Run(t, new(TodoE2ETestSuite))
}

func (s *TodoE2ETestSuite) TestTodoCRUD() {
    // Register test user
    registerReq := &helpers.RegisterRequest{
        Username: "testuser",
        Password: "testpass",
    }
    var registerRes helpers.RegisterResponse
    err := s.doRequest("POST", "/api/register", registerReq, &registerRes)
    s.Require().NoError(err, "Failed to register user")
    s.Require().NotEmpty(registerRes.ID, "Expected user ID to be returned")

    // Login
    loginReq := &helpers.LoginRequest{
        Username: "testuser",
        Password: "testpass",
    }
    var loginRes helpers.LoginResponse
    err = s.doRequest("POST", "/api/login", loginReq, &loginRes)
    s.Require().NoError(err, "Failed to login")
    s.Require().NotEmpty(loginRes.Token, "Expected token to be returned")
    s.token = loginRes.Token

    // Create todo
    createReq := &dto.CreateTodoRequest{
        Task:        "Test Todo",
        Description: "Test Description",
        Important:   true,
    }
    var createRes dto.CreateResponse
    err = s.doRequest("POST", "/api/todo", createReq, &createRes)
    s.Require().NoError(err, "Failed to create todo")
    s.Require().NotEmpty(createRes.ID, "Expected todo ID to be returned")

    // Get todos
    var todosRes dto.TodosResponse
    err = s.doRequest("GET", "/api/todos", nil, &todosRes)
    s.Require().NoError(err, "Failed to get todos")
    s.Require().Len(todosRes.Todos, 1, "Expected exactly one todo")
    s.Equal(createReq.Task, todosRes.Todos[0].Task)
    s.Equal(createReq.Description, todosRes.Todos[0].Description)
    s.Equal(createReq.Important, todosRes.Todos[0].Important)

    // Update todo
    updateReq := &dto.UpdateTodoRequest{
        ID:          createRes.ID,
        Task:        "Updated Todo",
        Description: "Updated Description",
        Important:   false,
        Done:        true,
    }
    var updateRes dto.SuccessResponse
    err = s.doRequest("PUT", "/api/todo/"+createRes.ID, updateReq, &updateRes)
    s.Require().NoError(err, "Failed to update todo")
    s.True(updateRes.Success)

    // Verify update
    err = s.doRequest("GET", "/api/todos", nil, &todosRes)
    s.Require().NoError(err, "Failed to get updated todo")
    s.Require().Len(todosRes.Todos, 1)
    s.Equal(updateReq.Task, todosRes.Todos[0].Task)
    s.Equal(updateReq.Description, todosRes.Todos[0].Description)
    s.Equal(updateReq.Important, todosRes.Todos[0].Important)
    s.Equal(updateReq.Done, todosRes.Todos[0].Done)

    // Delete todo
    err = s.doRequest("DELETE", "/api/todo/"+createRes.ID, nil, &updateRes)
    s.Require().NoError(err, "Failed to delete todo")
    s.True(updateRes.Success)

    // Verify deletion
    err = s.doRequest("GET", "/api/todos", nil, &todosRes)
    s.Require().NoError(err, "Failed to get todos after deletion")
    s.Empty(todosRes.Todos, "Expected no todos after deletion")
}

func (s *TodoE2ETestSuite) doRequest(method, path string, body interface{}, target interface{}) error {
    fullPath := path
    if s.token != "" {
        fullPath = path + "?token=" + s.token
    }
    return s.client.DoRequest(method, fullPath, body, target)
}