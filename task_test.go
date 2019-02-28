package main

import "testing"

var tests = []struct {
	title, content, expectedTitle, expectedContent string
}{
	{"Test", "This is a test", "Test", "This is a test"},
	{"0", "0", "0", "0"},
	{"", "", "", ""},
}

func TestNewTask(t *testing.T) {
	for _, test := range tests {
		if task := NewTask(test.title, test.content); task.Title != test.expectedTitle || task.Content != test.expectedContent || task.IsUrgent != false {
			t.Fatalf("NewTask(%s, %s) resolved in actual title: %s, actual content: %s, actual isUrgent status: %v, should be: %s, %s and %v.",
				test.title, test.content, task.Title, task.Content, task.IsUrgent, test.expectedTitle, test.expectedContent, false)
		}
	}
}
