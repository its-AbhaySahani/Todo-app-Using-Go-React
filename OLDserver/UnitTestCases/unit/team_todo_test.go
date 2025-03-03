package unit

import (
    "errors"
    "fmt"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
)

func TestCreateTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateTeamTodo")
    fmt.Println("Testing creating a new team todo")
    
    // Create a mock team todo repository
    mockRepo := new(mocks.MockTeamTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Setup expectations
    mockRepo.On("CreateTeamTodo", "Team Task", "Team Description", true, "team-1", "user-123").Return(
        models.TeamTodo{
            ID:          "team-todo-1",
            Task:        "Team Task",
            Description: "Team Description",
            Done:        false,
            Important:   true,
            TeamID:      "team-1",
            AssignedTo:  "user-123",
            Date:        dateStr,
            Time:        timeStr,
        }, nil)
    
    mockRepo.On("CreateTeamTodo", "Error Task", "Error Description", false, "invalid-team", "user-123").Return(
        models.TeamTodo{}, errors.New("team not found"))
    
    // Test successful team todo creation
    todo, err := mockRepo.CreateTeamTodo("Team Task", "Team Description", true, "team-1", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "team-todo-1", todo.ID)
    assert.Equal(t, "Team Task", todo.Task)
    assert.Equal(t, "Team Description", todo.Description)
    assert.False(t, todo.Done)
    assert.True(t, todo.Important)
    assert.Equal(t, "team-1", todo.TeamID)
    assert.Equal(t, "user-123", todo.AssignedTo)
    
    fmt.Printf("Created team todo with ID: %s\n", todo.ID)
    fmt.Printf("Team todo details: {Task:%s Description:%s Done:%t Important:%t TeamID:%s AssignedTo:%s}\n",
        todo.Task, todo.Description, todo.Done, todo.Important, todo.TeamID, todo.AssignedTo)
    
    // Test error case
    _, err = mockRepo.CreateTeamTodo("Error Task", "Error Description", false, "invalid-team", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("CreateTeamTodo test passed")
}

func TestGetTeamTodos(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTeamTodos")
    fmt.Println("Testing retrieving todos for a team")
    
    // Create a mock team todo repository
    mockRepo := new(mocks.MockTeamTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Create test data
    teamTodos := []models.TeamTodo{
        {
            ID:          "team-todo-1",
            Task:        "Team Task 1",
            Description: "Team Description 1",
            Done:        false,
            Important:   true,
            TeamID:      "team-1",
            AssignedTo:  "user-123",
            Date:        dateStr,
            Time:        timeStr,
        },
        {
            ID:          "team-todo-2",
            Task:        "Team Task 2",
            Description: "Team Description 2",
            Done:        true,
            Important:   false,
            TeamID:      "team-1",
            AssignedTo:  "user-456",
            Date:        dateStr,
            Time:        timeStr,
        },
    }
    
    // Setup expectations
    mockRepo.On("GetTeamTodos", "team-1").Return(teamTodos, nil)
    mockRepo.On("GetTeamTodos", "empty-team").Return([]models.TeamTodo{}, nil)
    mockRepo.On("GetTeamTodos", "invalid-team").Return([]models.TeamTodo{}, errors.New("team not found"))
    
    // Test getting todos for a team with todos
    todos, err := mockRepo.GetTeamTodos("team-1")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 2, len(todos))
    assert.Equal(t, "Team Task 1", todos[0].Task)
    assert.Equal(t, "Team Task 2", todos[1].Task)
    assert.Equal(t, "user-123", todos[0].AssignedTo)
    assert.Equal(t, "user-456", todos[1].AssignedTo)
    
    fmt.Printf("Retrieved %d team todos\n", len(todos))
    for i, todo := range todos {
        fmt.Printf("Team Todo %d: {ID:%s Task:%s AssignedTo:%s Done:%t Important:%t}\n",
            i+1, todo.ID, todo.Task, todo.AssignedTo, todo.Done, todo.Important)
    }
    
    // Test getting todos for a team with no todos
    emptyTodos, err := mockRepo.GetTeamTodos("empty-team")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyTodos))
    
    fmt.Println("Successfully retrieved empty todo list for team with no todos")
    
    // Test getting todos for an invalid team
    _, err = mockRepo.GetTeamTodos("invalid-team")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for invalid team: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTeamTodos test passed")
}

func TestUpdateTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateTeamTodo")
    fmt.Println("Testing updating a team todo")
    
    // Create a mock team todo repository
    mockRepo := new(mocks.MockTeamTodoRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    timeStr := currentTime.Format("15:04:05")
    
    // Setup expectations
    mockRepo.On("UpdateTeamTodo", "team-todo-1", "Updated Task", "Updated Description", true, false, "team-1", "user-456").Return(
        models.TeamTodo{
            ID:          "team-todo-1",
            Task:        "Updated Task",
            Description: "Updated Description",
            Done:        true,
            Important:   false,
            TeamID:      "team-1",
            AssignedTo:  "user-456",
            Date:        dateStr,
            Time:        timeStr,
        }, nil)
    
    mockRepo.On("UpdateTeamTodo", "nonexistent-todo", "Task", "Description", false, false, "team-1", "user-123").Return(
        models.TeamTodo{}, errors.New("team todo not found"))
    
    mockRepo.On("UpdateTeamTodo", "team-todo-2", "Task", "Description", false, false, "invalid-team", "user-123").Return(
        models.TeamTodo{}, errors.New("team not found"))
    
    // Test successful update
    todo, err := mockRepo.UpdateTeamTodo("team-todo-1", "Updated Task", "Updated Description", true, false, "team-1", "user-456")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "team-todo-1", todo.ID)
    assert.Equal(t, "Updated Task", todo.Task)
    assert.Equal(t, "Updated Description", todo.Description)
    assert.True(t, todo.Done)
    assert.False(t, todo.Important)
    assert.Equal(t, "user-456", todo.AssignedTo)
    
    fmt.Printf("Updated team todo: {ID:%s Task:%s Description:%s Done:%t Important:%t AssignedTo:%s}\n",
        todo.ID, todo.Task, todo.Description, todo.Done, todo.Important, todo.AssignedTo)
    
    // Test updating non-existent todo
    _, err = mockRepo.UpdateTeamTodo("nonexistent-todo", "Task", "Description", false, false, "team-1", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent team todo: %v\n", err)
    
    // Test updating todo for invalid team
    _, err = mockRepo.UpdateTeamTodo("team-todo-2", "Task", "Description", false, false, "invalid-team", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for invalid team: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("UpdateTeamTodo test passed")
}

func TestDeleteTeamTodo(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDeleteTeamTodo")
    fmt.Println("Testing deleting a team todo")
    
    // Create a mock team todo repository
    mockRepo := new(mocks.MockTeamTodoRepository)
    
    // Setup expectations
    mockRepo.On("DeleteTeamTodo", "team-todo-1", "team-1").Return(nil)
    mockRepo.On("DeleteTeamTodo", "nonexistent-todo", "team-1").Return(errors.New("team todo not found"))
    mockRepo.On("DeleteTeamTodo", "team-todo-2", "invalid-team").Return(errors.New("team not found"))
    
    // Test successful deletion
    err := mockRepo.DeleteTeamTodo("team-todo-1", "team-1")
    assert.NoError(t, err)
    
    fmt.Println("Team todo successfully deleted")
    
    // Test deleting non-existent todo
    err = mockRepo.DeleteTeamTodo("nonexistent-todo", "team-1")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent team todo: %v\n", err)
    
    // Test deleting todo from invalid team
    err = mockRepo.DeleteTeamTodo("team-todo-2", "invalid-team")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "team not found")
    
    fmt.Printf("Correctly got error for invalid team: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("DeleteTeamTodo test passed")
}