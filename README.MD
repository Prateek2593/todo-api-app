# To-Do List REST API

A simple, reliable RESTful API for managing a to-do list, built with Go and `gorilla/mux`. Create, read, update, and delete tasks with JSON persistence (`todos.json`). Features include input validation (non-empty titles, valid priorities, UUIDs) and a clean, extensible design inspired by a CLI to-do app.

Perfect for learning Go-based APIs, testing with Postman, or integrating with a front-end.

## Features

- **CRUD Operations**:
  - Create tasks (`POST /todos`)
  - List all tasks (`GET /todos`)
  - Fetch a task by ID (`GET /todos/{id}`)
  - Update tasks (`PUT /todos/{id}`)
  - Delete tasks (`DELETE /todos/{id}`)
- **Validation**:
  - Non-empty `title` for tasks.
  - `priority` must be `low`, `medium`, or `high` (optional).
  - Task IDs must be valid UUIDs.
- **Persistence**: Stores tasks in `todos.json`.
- **Error Handling**: Clear JSON errors (e.g., `{"error": "Title is required"}`).

## Prerequisites

- [Go](https://golang.org/dl/) (1.20 or later)
- [Postman](https://www.postman.com/downloads/) (optional, for testing)
- Basic knowledge of HTTP and JSON

## Installation

1. **Clone the Repository**:
   ```bash
   git clone github.com/Prateek2593/todo-api-app
   cd todo-api-app
   ```
2. **Install Dependencies:**
   ```bash
   go mod tidy
   ```
3. **Run the API:**
   ```bash
   go run .
   ```

## Usage

Send HTTP requests to `http://localhost:8080` using Postman, curl, or a client. The API stores tasks in `todos.json` and validates inputs.

## Endpoints

| Method | Endpoint      | Description       | Request Body Example                          | Response Example                                         |
| ------ | ------------- | ----------------- | --------------------------------------------- | -------------------------------------------------------- |
| POST   | `/todos`      | Create a new task | `{"title": "Write Code", "priority": "high"}` | `201: {"id": "123e4567...", "title": "Write Code", ...}` |
| GET    | `/todos`      | List all tasks    | None                                          | `200: [{"id": "123e4567...", ...}, ...]`                 |
| GET    | `/todos/{id}` | Get a task by ID  | None                                          | `200: {"id": "123e4567...", ...}` or `404`               |
| PUT    | `/todos/{id}` | Update a task     | `{"completed": true, "priority": "medium"}`   | `200: {"id": "123e4567...", ...}` or `404`               |
| DELETE | `/todos/{id}` | Delete a task     | None                                          | `204 (no content)` or `404`                              |

## Project Structure

- **`main.go`**: Entry point.
- **`todo.go`**: Defines `Todo` struct and methods.
- **`storage.go`**: JSON file storage logic
- **`todos.json``**: JSON file for storing tasks.
- **`go.mod`**: Go module dependencies.

## Notes

- **Persistence**: Tasks are saved to `todos.json`. Delete it to reset the list.
- **Validation**: Ensures non-empty titles, valid priorities (low, medium, high), and UUID IDs.
- **Extensibility**: Ready for features like sorting (`?sort=createdAt`), pagination, or a front-end.
- **Testing**: Use Postman for interactive testing. Set `Content-Type: application/json` for POST and PUT.

## Future Improvements

- Add sorting (`GET /todos?sort=createdAt`).
- Implement pagination (`GET /todos?page=1&limit=10`).
- Support due dates in `Todo` struct.
- Integrate a database (e.g., SQLite) for scalability.
