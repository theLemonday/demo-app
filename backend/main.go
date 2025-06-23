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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	todos         []Todo
	mu            sync.Mutex
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"method", "path"},
	)

	todosCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "todos_created_total",
			Help: "Number of todos created",
		},
	)
)

func init() {
	prometheus.MustRegister(totalRequests, todosCreated)
}

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

func deleteTodo(id int) {
	mu.Lock()
	defer mu.Unlock()
	var newTodos []Todo
	for _, t := range todos {
		if id != t.ID {
			newTodos = append(newTodos, t)
		}
	}
	todos = newTodos
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	totalRequests.WithLabelValues("GET", "/api/todos").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getAllTodos())
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	totalRequests.WithLabelValues("POST", "/api/todos").Inc()

	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newTodo.Title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	created := createTodo(newTodo)

	todosCreated.Inc()

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

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	totalRequests.WithLabelValues("DELETE", "/api/todos/{id}").Inc()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err)
	}
	deleteTodo(id)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
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

	r.With(basicAuth).Get("/api/todos", getTodosHandler)
	r.With(basicAuth, authorizeAdmin).Post("/api/todos", addTodoHandler)
	r.With(basicAuth, authorizeAdmin).Post("/api/todos/{id}/toggle", toggleTodoHandler)
	r.With(basicAuth, authorizeAdmin).Delete("/api/todos/{id}", deleteTodoHandler)

	r.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
