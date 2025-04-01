package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func loadTasks() []Task {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}
		}
		log.Fatal(err)
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Fatal(err)
	}
	return tasks
}

func saveTasks(tasks []Task) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("tasks.json", data, 0644); err != nil {
		log.Fatal(err)
	}
}

func addTask(description string) {
	if strings.TrimSpace(description) == "" {
		log.Fatal("Error: Description cannot be empty.")
	}
	tasks := loadTasks()
	newTask := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("Task added successfully (ID: %d)\n", len(tasks))
}

func updateTask(id int, newDescription string) {
	if strings.TrimSpace(newDescription) == "" {
		log.Fatal("Error: Description cannot be empty.")
	}
	tasks := loadTasks()
	flag := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			flag = true
			break
		}
	}
	if !flag {
		log.Fatalf("Error: Task with ID %d not found.", id)
	}
	saveTasks(tasks)
	fmt.Printf("Task %d updated successfully.\n", id)
}

func deleteTask(id int) {
	tasks := loadTasks()
	newTasks := []Task{}
	flag := false
	for _, task := range tasks {
		if task.ID == id {
			flag = true
		} else {
			newTasks = append(newTasks, task)
		}
	}
	if !flag {
		log.Fatalf("Error: Task with ID %d not found.", id)
	}
	saveTasks(newTasks)
	fmt.Printf("Task %d deleted successfully.\n", id)
}

func markTaskStatus(id int, status string) {
	tasks := loadTasks()
	flag := false
	for i, task := range tasks {
		if task.ID == id {
			if task.Status == status {
				fmt.Printf("Task %d is already %s.\n", id, status)
				return
			}
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			flag = true
			break
		}
	}
	if !flag {
		log.Fatalf("Error: Task with ID %d not found.", id)
	}
	saveTasks(tasks)
	fmt.Printf("Task %d marked as %s.\n", id, status)
}

func listTasks(status string) {
	tasks := loadTasks()
	filtred := []Task{}
	for _, task := range tasks {
		if task.Status == status || status == "" {
			filtred = append(filtred, task)
		}
	}
	if len(filtred) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	for _, task := range filtred {
		fmt.Printf("Task %d:\n", task.ID)
		fmt.Printf("  Description: %s\n", task.Description)
		fmt.Printf("  Status: %s\n", task.Status)
		fmt.Printf("  Created: %s\n", task.CreatedAt.Format(time.RFC3339))
		fmt.Printf("  Updated: %s\n", task.UpdatedAt.Format(time.RFC3339))
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli [command]")
		fmt.Println("Commands: add, update, delete, mark-in-progress, mark-done, list")
		os.Exit(1)
	}
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("Error: Missing description for add command.")
		}
		addTask(os.Args[2])
	case "update":
		if len(os.Args) < 4 {
			log.Fatal("Error: Missing arguments for update command.")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid task ID.")
		}
		updateTask(id, os.Args[3])
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Error: Missing task ID for delete command.")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid task ID")
		}
		deleteTask(id)
	case "mark-in-progress":
		if len(os.Args) < 3 {
			log.Fatal("Error: Missing task ID for mark-in-progress command.")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid task ID.")
		}
		markTaskStatus(id, "in-progress")
	case "mark-done":
		if len(os.Args) < 3 {
			log.Fatal("Error: Missing task ID for mark-done command.")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error: Invalid task ID.")
		}
		markTaskStatus(id, "done")
	case "list":
		statusFilt := ""
		if len(os.Args) >= 3 {
			statusFilt := os.Args[2]
			valid := map[string]bool{
				"todo":        true,
				"in-progress": true,
				"done":        true,
			}
			if !valid[statusFilt] {
				log.Fatal("Error: Invalid status. Use 'todo', 'in-progress', or 'done'.")
			}
		}
		listTasks(statusFilt)
	default:
		fmt.Printf("Error: Unknown command '%s'.\n", command)
		fmt.Println("Available commands: add, update, delete, mark-in-progress, mark-done, list")
		os.Exit(1)
	}
}
