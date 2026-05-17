package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// TaskManager acts as our core service engine
type TaskManager struct {
	filename string
	tasks    []Task
}

// NewTaskManager is a "constructor" function helper
func NewTaskManager(filename string) *TaskManager {
	return &TaskManager{
		filename: filename,
		tasks:    []Task{},
	}
}

// Save persists the internal tasks slice back to the file
func (m *TaskManager) Save() error {
	jsonTasks, err := json.MarshalIndent(m.tasks, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filename, jsonTasks, 0644)
}

func renderCheckbox(b bool) string {
	if b {
		return "[x]"
	}
	return "[ ]"
}

func (m *TaskManager) Load() error {
	data, err := os.ReadFile(m.filename)
	if err == nil {
		json.Unmarshal(data, &m.tasks)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal("Could not read file:", err)
	}
	return err
}

func (m *TaskManager) List() {
	if len(m.tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	for _, t := range m.tasks {
		fmt.Printf("%d. %s %s\n", t.ID, renderCheckbox(t.Done), t.Title)
	}
}

func (m *TaskManager) Add(title string) {
	newID := 1
	if len(m.tasks) > 0 {
		newID = m.tasks[len(m.tasks)-1].ID + 1
	}
	m.tasks = append(m.tasks, Task{ID: newID, Title: title, Done: false})
}

// Done looks for the task ID, marks it done, and returns true if found.
func (m *TaskManager) Done(id int) bool {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks[i].Done = true
			return true
		}
	}
	return false
}

// Delete removes the task by ID and returns true if found.
func (m *TaskManager) Delete(id int) bool {
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			m.tasks = slices.Delete(m.tasks, i, i+1)
			return true
		}
	}
	return false
}

func main() {
	// Initialize our engine
	manager := NewTaskManager("tasks.json")

	// 1. LOAD DATA
	if err := manager.Load(); err != nil {
		log.Fatal("Failed to load tasks:", err)
	}

	// 2. PARSE CLI INPUT
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-manager [add|list|done|delete]")
		return
	}

	command := os.Args[1]

	// 3. EXECUTE ENGINE ACTIONS
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task title")
			return
		}
		manager.Add(os.Args[2])
		fmt.Println("Task added successfully!")

	case "list":
		manager.List()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid ID")
			return
		}

		if manager.Done(id) {
			fmt.Println("Task done!")
		} else {
			fmt.Println("Task not found.")
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("Error: Invalid task ID. Please provide a valid number.")
			return
		}

		if manager.Delete(id) {
			fmt.Printf("Task %d deleted!\n", id)
		} else {
			fmt.Println("Task not found.")
		}

	default:
		fmt.Println("Unknown command:", command)
		return
	}

	// 4. SAVE DATA
	if err := manager.Save(); err != nil {
		log.Fatal("Failed to save tasks:", err)
	}
}

/*
func main() {
	filename := "tasks.json"
	var tasks []Task

	// 1. LOAD: Read existing tasks
	data, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(data, &tasks)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal("Could not read file:", err)
	}

	// 2. CHECK: Ensure we have a command
	if len(os.Args) < 2 {
		fmt.Println("Usage: taskmaster [add|list|done]")
		return
	}

	command := os.Args[1]

	// 3. ACT: Modify the data in memory
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task title")
			return
		}
		newID := 1
		if len(tasks) > 0 {
			newID = tasks[len(tasks)-1].ID + 1
		}
		tasks = append(tasks, Task{ID: newID, Title: os.Args[2], Done: false})
		fmt.Println("Task added!")

	case "list":
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}
		for _, t := range tasks {
			fmt.Printf("%d. %s %s\n", t.ID, renderCheckbox(t.Done), t.Title)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task ID")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				found = true
				break
			}
		}
		if found {
			fmt.Printf("Task %d marked as done!\n", id)
		} else {
			fmt.Println("Task not found.")
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Provide a task ID")
			return
		}

		id, _ := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID. Please provide a valid number.")
			return
		}

		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks = slices.Delete(tasks, i, i+1)
				found = true
				break
			}
		}

		if found {
			fmt.Printf("Task %d deleted!\n", id)
		} else {
			fmt.Println("Task not found.")
		}

	default:
		fmt.Println("Unknown command:", command)
		return
	}

	// 4. SAVE: Write changes back to disk
	jsonTasks, _ := json.MarshalIndent(tasks, "", "	")
	os.WriteFile(filename, jsonTasks, 0644)
}
*/
