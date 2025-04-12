package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func main() {
	// initialize a new Todos instance and a Storage instance for the "todos.json" file
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	if err := storage.Load(&todos); err != nil {
		log.Fatalf("Failed to load todos: %v", err)
	}

	// route handler for /todos
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// handle GET request to retrieve todos
			listTodos(w, r, &todos)
		case http.MethodPost:
			// handle POST request to add a new todo
			addTodo(w, r, &todos, storage)
		default:
			// handle unsupported methods
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			id := r.URL.Path[len("/todos/"):]
			deleteTodo(w, r, &todos, storage, id)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// listTodos handles the GET request to retrieve todos and send them as a JSON response
func listTodos(w http.ResponseWriter, r *http.Request, todos *Todos) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
	}
}

// addTodo handles the POST request to add a new todo and save it to the file
func addTodo(w http.ResponseWriter, r *http.Request, todos *Todos, storage *Storage[Todos]) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode todo", http.StatusBadRequest)
		return
	}
	todo.ID = uuid.NewString() // generate a new UUID for the todo
	todo.CreatedAt = time.Now()
	*todos = append(*todos, todo)
	if err := storage.Save(*todos); err != nil { // save the updated todos to the file
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
func deleteTodo(w http.ResponseWriter, r *http.Request, todos *Todos, storage *Storage[Todos], id string) {
	for i, todo := range *todos {
		if todo.ID == id {
			*todos = append((*todos)[:i], (*todos)[i+1:]...)
			if err := storage.Save(*todos); err != nil {
				http.Error(w, "Failed to save todos", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent) // send a 204 No Content response
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}
