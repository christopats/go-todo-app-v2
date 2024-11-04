package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", showTodoList)
	http.HandleFunc("/add", addTodo)
	http.Handle("/static/", http.FileServer(http.Dir("public")))

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// A Todo represents a task in our todo list
type Todo struct {
	ID     int
	Task   string
	Status bool
}

// This slice will act as our in-memory database for now
var todoList = []Todo{
	{ID: 1, Task: "Learn Go", Status: false},
	{ID: 2, Task: "Write a blog post", Status: false},
}

func showTodoList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, todoList)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	// Add the new task to our in-memory database
	newID := len(todoList) + 1
	todoList = append(todoList, Todo{ID: newID, Task: task, Status: false})

	// Rerender the todo list template
	tmpl, err := template.ParseFiles("templates/todo-list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, todoList)
}

