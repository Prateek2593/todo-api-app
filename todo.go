package main

import "time"

type Todo struct {
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Priority    string     `json:"priority"`
	Notes       string     `json:"notes,omitempty"`
}

type Todos []Todo
