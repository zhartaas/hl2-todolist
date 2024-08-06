package main

import (
	"sync"
	"time"
)

var LocalMap sync.Map

type TaskInput struct {
	Title    string `json:"title" example:"Прочитать книгу" extensions:"x-order=1"`
	ActiveAt string `json:"activeAt" example:"2026-01-02" extensions:"x-order=2"`
}

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Done     bool   `json:"-"`
}

type TaskResponse struct {
	ID string `json:"id"`
}

type GetTasksResponse []*Task

func (s GetTasksResponse) Len() int {
	return len(s)
}
func (s GetTasksResponse) Less(i, j int) bool {
	datei, _ := time.Parse(time.DateOnly, s[i].ActiveAt)
	datej, _ := time.Parse(time.DateOnly, s[j].ActiveAt)

	return datej.After(datei)
}
func (s GetTasksResponse) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
