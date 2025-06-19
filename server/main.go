package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	todos []Todo
	mu    sync.Mutex
)

func getAllTodos() []Todo {
	mu.Lock()
	defer mu.Unlock()
	if len(todos) == 0 {
		return []Todo{}
	}
	return slices.Clone(todos)
}

func createTodo(newTodo Todo) Todo {
	mu.Lock()
	defer mu.Unlock()
	newTodo.ID = rand.Intn(1000000)
	todos = append(todos, newTodo)
	return newTodo
}

func toggleTodoByID(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = !todos[i].Done
			return true
		}
	}
	return false
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getAllTodos())
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created := createTodo(newTodo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func toggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if !toggleTodoByID(id) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)

	r.Get("/api/todos", getTodosHandler)
	r.Post("/api/todos", addTodoHandler)
	r.Post("/api/todos/{id}/toggle", toggleTodoHandler)

	log.Printf("Server running on http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
