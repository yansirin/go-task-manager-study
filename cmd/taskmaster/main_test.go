package main

import "testing"

// MockStore satisfies the TaskStore interface without touching the disk
type MockStore struct {
	tasks []Task
}

func (m *MockStore) Load() ([]Task, error) {
	return m.tasks, nil
}

func (m *MockStore) Save(tasks []Task) error {
	m.tasks = tasks
	return nil
}

func TestTaskManager_Add(t *testing.T) {
	// Setup
	mock := &MockStore{tasks: []Task{}}
	manager := NewTaskManager(mock)
	manager.Load()

	// Act
	manager.Add("Learn Go Testing")

	// Assert
	if len(manager.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(manager.tasks))
	}
	if manager.tasks[0].Title != "Learn Go Testing" {
		t.Errorf("Expected title 'Learn Go Testing', got %s", manager.tasks[0].Title)
	}
}

func TestTaskManager_Delete(t *testing.T) {
	// Arrange: Create a fake world with two hardcoded tasks
	initialTasks := []Task{
		{ID: 1, Title: "Buy Milk", Done: false},
		{ID: 2, Title: "Learn Go Testing", Done: false},
	}

	mock := &MockStore{tasks: initialTasks}
	manager := NewTaskManager(mock)

	// Pull the hardcoded tasks out of the store and into the manager's memory
	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	// Act: Execute the action we are trying to test
	wasDeleted := manager.Delete(1)

	// Assert: Verify the results match our expectations
	if !wasDeleted {
		t.Errorf("Expected Delete(1) to return true, but it returned false")
	}

	// Verify that the task LEFT behind is indeed Task ID 2
	if manager.tasks[0].ID != 2 {
		t.Errorf("Expected remaining task to have ID 2, but got ID %d", manager.tasks[0].ID)
	}

}

func TestTaskManager_Done(t *testing.T) {
	// Arrange
	initialTasks := []Task{
		{ID: 1, Title: "Buy Milk", Done: false},
		{ID: 2, Title: "Learn Go Testing", Done: false},
	}

	mock := &MockStore{tasks: initialTasks}
	manager := NewTaskManager(mock)

	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	// Act
	wasDone := manager.Done(1)

	// Assert
	if !wasDone {
		t.Errorf("Expected Done(1) to return true, but it returned false")
	}

	if manager.tasks[0].Done != true {
		t.Errorf("Expected tested task to have Done true, but got Done false")
	}
}
