package main

import "time"

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Priority    string     `json:"priority"`
	Notes       string     `json:"notes"`
}

type Todos []Todo
