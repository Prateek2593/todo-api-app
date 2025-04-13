package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// App represents the main application
type App struct {
	todos   *Todos
	storage *Storage[Todos]
}

func main() {
	// initialize a new Todos instance and a Storage instance for the "todos.json" file
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	if err := storage.Load(&todos); err != nil {
		log.Fatalf("Failed to load todos: %v", err)
	}

	app := &App{
		todos:   &todos,
		storage: storage,
	}

	// create a new router using gorilla/mux
	router := mux.NewRouter()
	// define the routes and their handlers
	router.HandleFunc("/todos", app.listTodos).Methods("GET")
	router.HandleFunc("/todos", app.addTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", app.getTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", app.updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", app.deleteTodo).Methods("DELETE")

	// start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// listTodos handles the GET request to retrieve todos and send them as a JSON response
func (app *App) listTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(*app.todos); err != nil {
		http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
	}
}

// addTodo handles the POST request to add a new todo and save it to the file
func (app *App) addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode todo", http.StatusBadRequest)
		return
	}

	// validate the todo title
	if strings.TrimSpace(todo.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// validate the todo priority
	if todo.Priority != "" {
		switch strings.ToLower(todo.Priority) {
		case "low", "medium", "high":
			// valid priority
			todo.Priority = strings.ToLower(todo.Priority) // normalize the priority to lowercase
		default:
			http.Error(w, "Invalid priority. Allowed values are: low, medium, high", http.StatusBadRequest)
			return

		}
	}

	todo.ID = uuid.NewString() // generate a new UUID for the todo
	todo.CreatedAt = time.Now()
	*app.todos = append(*app.todos, todo)
	if err := app.storage.Save(*app.todos); err != nil { // save the updated todos to the file
		http.Error(w, "Failed to save todos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // send a 201 Created response
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}

// deleteTodo handles the DELETE request to delete a todo by its ID and save the updated list to the file
func (app *App) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	for i, todo := range *app.todos {
		if todo.ID == id {
			*app.todos = append((*app.todos)[:i], (*app.todos)[i+1:]...)
			if err := app.storage.Save(*app.todos); err != nil {
				http.Error(w, "Failed to save todos", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent) // send a 204 No Content response
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

// updateTodo handles the PUT request to update a todo by its ID and save the updated list to the file
func (app *App) updateTodo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	var updates struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
		Priority  *string `json:"priority"`
		Notes     *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Failed to decode updates", http.StatusBadRequest)
		return
	}

	// validate the updates
	if updates.Title == nil && updates.Completed == nil && updates.Priority == nil && updates.Notes == nil {
		http.Error(w, "At least one field must be updated", http.StatusBadRequest)
		return
	}

	if updates.Title != nil && strings.TrimSpace(*updates.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if updates.Priority != nil {
		switch strings.ToLower(*updates.Priority) {
		case "low", "medium", "high":
			// valid priority
			*updates.Priority = strings.ToLower(*updates.Priority) // normalize the priority to lowercase
		default:
			http.Error(w, "Invalid priority. Allowed values are: low, medium, high", http.StatusBadRequest)
			return
		}
	}

	for i, todo := range *app.todos {
		if todo.ID == id {
			if updates.Title != nil {
				todo.Title = *updates.Title
			}
			if updates.Completed != nil {
				todo.Completed = *updates.Completed
				if todo.Completed {
					now := time.Now()
					todo.CompletedAt = &now
				} else {
					todo.CompletedAt = nil
				}
			}
			if updates.Priority != nil {
				todo.Priority = *updates.Priority
			}
			if updates.Notes != nil {
				todo.Notes = *updates.Notes
			}
			(*app.todos)[i] = todo                               // update the todo in the slice
			if err := app.storage.Save(*app.todos); err != nil { // save the updated todos to the file
				http.Error(w, "Failed to save todos", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(todo); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

// getTodo handles the GET request to retrieve a specific todo by its ID and send it as a JSON response
func (app *App) getTodo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	for _, todo := range *app.todos {
		if todo.ID == id {
			if err := json.NewEncoder(w).Encode(todo); err != nil {
				http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
				return
			}
			return
		}

	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}
