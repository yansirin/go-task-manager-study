package main

import (
	"encoding/json"
	"errors"
	"flag"
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

type TaskStore interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}

type JSONFileStore struct {
	filename string
}

func NewJSONFileStore(filename string) *JSONFileStore {
	return &JSONFileStore{filename: filename}
}

// TaskManager acts as our core service engine
type TaskManager struct {
	store TaskStore
	tasks []Task
}

// NewTaskManager is a "constructor" function helper
func NewTaskManager(store TaskStore) *TaskManager {
	return &TaskManager{
		store: store,
		tasks: []Task{},
	}
}

func renderCheckbox(b bool) string {
	if b {
		return "[x]"
	}
	return "[ ]"
}

func (s *JSONFileStore) Load() ([]Task, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *JSONFileStore) Save(tasks []Task) error {
	jsonTasks, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, jsonTasks, 0644)
}

func (m *TaskManager) Load() error {
	loadedTasks, err := m.store.Load()
	if err != nil {
		return err
	}
	m.tasks = loadedTasks
	return nil
}

func (m *TaskManager) Save() error {
	return m.store.Save(m.tasks)
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
	fileStore := NewJSONFileStore("tasks.json")

	// Initialize our engine
	manager := NewTaskManager(fileStore)

	// 1. LOAD DATA
	if err := manager.Load(); err != nil {
		log.Fatal("Failed to load tasks:", err)
	}

	// 2. PARSE CLI INPUT
	if len(os.Args) < 2 {
		fmt.Println("Usage: taskmaster <command> [<args>]")
		fmt.Println("Commands: add, list, done, delete")
		return
	}

	// 3. DEFINE SUBCOMMANDS
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// 4. EXECUTE ENGINE ACTIONS
	switch os.Args[1] {
	case "add":
		// Parse everything after 'add'
		addCmd.Parse(os.Args[2:])

		// flag.Args() gathers everything remaining after flags are parsed
		remaining := addCmd.Args()
		if len(remaining) < 1 {
			fmt.Println("Error: Provide a task title")
			return
		}

		manager.Add(remaining[0])
		fmt.Println("Task added successfully!")

	case "list":
		manager.List()

	case "done":
		doneCmd.Parse(os.Args[2:])

		remaining := doneCmd.Args()
		if len(remaining) < 1 {
			fmt.Println("Error: Provide a task ID")
			return
		}
		id, err := strconv.Atoi(remaining[0])
		if err != nil {
			fmt.Println("Error: Invalid ID")
			return
		}

		if manager.Done(id) {
			fmt.Printf("Task %d marked as done!\n", id)
		} else {
			fmt.Println("Task not found.")
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])

		remaining := deleteCmd.Args()
		if len(remaining) < 1 {
			fmt.Println("Error: Provide a task ID")
			return
		}
		id, err := strconv.Atoi(remaining[0])

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
		fmt.Println("Unknown command:", os.Args[1])
		return
	}

	// 4. SAVE DATA
	if err := manager.Save(); err != nil {
		log.Fatal("Failed to save tasks:", err)
	}
}
