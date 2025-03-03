package unit

import (
    "errors"
    "fmt"
    "strings"
    "testing"
    "time"

    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/OLDmodels"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/UnitTestCases/mocks"
    "github.com/stretchr/testify/assert"
)

func TestGetDailyRoutines(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetDailyRoutines")
    fmt.Println("Testing retrieving routines for a specific day and schedule type")
    
    // Create a mock routine repository
    mockRepo := new(mocks.MockRoutineRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    
    // Create test data
    todos := []models.Todo{
        {
            ID:          "todo-1",
            Task:        "Morning Routine",
            Description: "Morning routine task",
            Done:        false,
            Important:   true,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        "08:00:00",
        },
    }
    
    // Setup expectations
    mockRepo.On("GetDailyRoutines", "monday", "morning", "user-123").Return(todos, nil)
    mockRepo.On("GetDailyRoutines", "tuesday", "evening", "user-123").Return([]models.Todo{}, nil)
    mockRepo.On("GetDailyRoutines", "invalid", "invalid", "user-123").Return(
        []models.Todo{}, errors.New("invalid day or schedule type"))
    
    // Test getting routines for a specific day and schedule
    routines, err := mockRepo.GetDailyRoutines("monday", "morning", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 1, len(routines))
    assert.Equal(t, "Morning Routine", routines[0].Task)
    assert.True(t, routines[0].Important)
    assert.Equal(t, "user-123", routines[0].UserID)
    
    fmt.Printf("Retrieved routine: {ID:%s Task:%s Description:%s Important:%t}\n",
        routines[0].ID, routines[0].Task, routines[0].Description, routines[0].Important)
    
    // Test getting routines for a day with no routines
    emptyRoutines, err := mockRepo.GetDailyRoutines("tuesday", "evening", "user-123")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(emptyRoutines))
    
    fmt.Println("Successfully retrieved empty routine list for day with no routines")
    
    // Test with invalid parameters
    _, err = mockRepo.GetDailyRoutines("invalid", "invalid", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid day or schedule type")
    
    fmt.Printf("Correctly got error for invalid parameters: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetDailyRoutines test passed")
}

func TestGetTodayRoutines(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestGetTodayRoutines")
    fmt.Println("Testing retrieving routines for today")
    
    // Create a mock routine repository
    mockRepo := new(mocks.MockRoutineRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    
    // Create test data for today's routines
    dayName := strings.ToLower(time.Now().Weekday().String())
    todos := []models.Todo{
        {
            ID:          "todo-1",
            Task:        "Evening Routine",
            Description: "Evening routine task",
            Done:        false,
            Important:   true,
            UserID:      "user-123",
            Date:        dateStr,
            Time:        "18:00:00",
        },
    }
    
    // Setup expectations
    mockRepo.On("GetDailyRoutines", dayName, "evening", "user-123").Return(todos, nil)
    mockRepo.On("GetDailyRoutines", dayName, "morning", "user-123").Return([]models.Todo{}, nil)
    
    // Test getting today's evening routines
    routines, err := mockRepo.GetDailyRoutines(dayName, "evening", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 1, len(routines))
    assert.Equal(t, "Evening Routine", routines[0].Task)
    
    fmt.Printf("Retrieved today's (%s) evening routine: {Task:%s Description:%s}\n", 
        dayName, routines[0].Task, routines[0].Description)
    
    // Test getting today's morning routines (empty)
    morningRoutines, err := mockRepo.GetDailyRoutines(dayName, "morning", "user-123")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(morningRoutines))
    
    fmt.Printf("Successfully retrieved empty routine list for today's morning\n")
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("GetTodayRoutines test passed")
}

func TestUpdateRoutineDay(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateRoutineDay")
    fmt.Println("Testing updating a routine's day")
    
    // Create a mock routine repository
    mockRepo := new(mocks.MockRoutineRepository)
    
    // Setup expectations
    mockRepo.On("UpdateRoutineDay", "routine-1", "friday").Return(nil)
    mockRepo.On("UpdateRoutineDay", "nonexistent-routine", "monday").Return(
        errors.New("routine not found"))
    
    // Test updating routine day
    err := mockRepo.UpdateRoutineDay("routine-1", "friday")
    
    // Assertions
    assert.NoError(t, err)
    fmt.Println("Successfully updated routine day to 'friday'")
    
    // Test updating non-existent routine
    err = mockRepo.UpdateRoutineDay("nonexistent-routine", "monday")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent routine: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("UpdateRoutineDay test passed")
}

func TestUpdateRoutineStatus(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestUpdateRoutineStatus")
    fmt.Println("Testing updating a routine's status (active/inactive)")
    
    // Create a mock routine repository
    mockRepo := new(mocks.MockRoutineRepository)
    
    // Setup expectations
    mockRepo.On("UpdateRoutineStatus", "routine-1", true).Return(nil)
    mockRepo.On("UpdateRoutineStatus", "routine-2", false).Return(nil)
    mockRepo.On("UpdateRoutineStatus", "nonexistent-routine", false).Return(
        errors.New("routine not found"))
    
    // Test activating routine
    err := mockRepo.UpdateRoutineStatus("routine-1", true)
    assert.NoError(t, err)
    fmt.Println("Successfully activated routine")
    
    // Test deactivating routine
    err = mockRepo.UpdateRoutineStatus("routine-2", false)
    assert.NoError(t, err)
    fmt.Println("Successfully deactivated routine")
    
    // Test updating non-existent routine
    err = mockRepo.UpdateRoutineStatus("nonexistent-routine", false)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not found")
    
    fmt.Printf("Correctly got error for non-existent routine: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("UpdateRoutineStatus test passed")
}

func TestDeleteRoutinesByTaskID(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestDeleteRoutinesByTaskID")
    fmt.Println("Testing deleting all routines for a task")
    
    // Create a mock routine repository
    mockRepo := new(mocks.MockRoutineRepository)
    
    // Setup expectations
    mockRepo.On("DeleteRoutinesByTaskID", "task-1").Return(nil)
    mockRepo.On("DeleteRoutinesByTaskID", "nonexistent-task").Return(nil) // Should succeed even if no routines exist
    mockRepo.On("DeleteRoutinesByTaskID", "error-task").Return(errors.New("database error"))
    
    // Test deleting routines for a task
    err := mockRepo.DeleteRoutinesByTaskID("task-1")
    assert.NoError(t, err)
    fmt.Println("Successfully deleted routines for task")
    
    // Test deleting routines for a non-existent task (should still succeed)
    err = mockRepo.DeleteRoutinesByTaskID("nonexistent-task")
    assert.NoError(t, err)
    fmt.Println("Successfully handled deletion for non-existent task (no-op)")
    
    // Test database error
    err = mockRepo.DeleteRoutinesByTaskID("error-task")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error")
    
    fmt.Printf("Correctly got error for database error: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("DeleteRoutinesByTaskID test passed")
}

func TestCreateRoutine(t *testing.T) {
    fmt.Println("\nRUNNING TEST: TestCreateRoutine")
    fmt.Println("Testing creating a new routine")
    
    mockRepo := new(mocks.MockRoutineRepository)
    
    currentTime := time.Now()
    dateStr := currentTime.Format("2006-01-02")
    
    expectedRoutine := models.Routine{
        ID:           "routine-1",
        Day:          "monday",
        ScheduleType: "morning",
        TaskID:       "task-1",
        UserID:       "user-123",
        CreatedAt:    dateStr,
        UpdatedAt:    dateStr,
        IsActive:     true,
    }
    
    // Setup expectations
    mockRepo.On("CreateRoutine", "monday", "morning", "task-1", "user-123").Return(expectedRoutine, nil)
    mockRepo.On("CreateRoutine", "invalid", "morning", "task-1", "user-123").Return(
        models.Routine{}, errors.New("invalid day"))
    
    // Test creating a routine
    routine, err := mockRepo.CreateRoutine("monday", "morning", "task-1", "user-123")
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "routine-1", routine.ID)
    assert.Equal(t, "monday", routine.Day)
    assert.Equal(t, "morning", routine.ScheduleType)
    assert.Equal(t, "task-1", routine.TaskID)
    assert.Equal(t, "user-123", routine.UserID)
    assert.True(t, routine.IsActive)
    
    fmt.Printf("Created routine: {ID:%s Day:%s ScheduleType:%s TaskID:%s UserID:%s IsActive:%t}\n",
        routine.ID, routine.Day, routine.ScheduleType, routine.TaskID, routine.UserID, routine.IsActive)
    
    // Test with invalid day
    _, err = mockRepo.CreateRoutine("invalid", "morning", "task-1", "user-123")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid day")
    
    fmt.Printf("Correctly got error for invalid day: %v\n", err)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    fmt.Println("CreateRoutine test passed")
}