package main

// Task is the main struct of this Todo application
type Task struct {
	ID       int
	Title    string
	Content  string
	IsUrgent bool
}

// NewTask is the contructor of the task struct
func NewTask(title string, content string) *Task {
	task := Task{
		Title:   title,
		Content: content,
	}
	return &task
}
