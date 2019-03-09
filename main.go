package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	serverPort = ":3000"
	layout     = "templates/layouts/layout.html"
)

var baseHTML string

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/submitTask", submitTaskHandler)
	http.HandleFunc("/deleteTask", deleteTaskHandler)
	http.HandleFunc("/task", singleTaskHandler)
	http.HandleFunc("/editTask", editTaskHandler)
	log.Printf("GOTTA-DO application is up and running on port %v.", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	q := `SELECT * FROM tasks`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("Error while retrieving tasks: %v", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.IsUrgent)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		tasks = append(tasks, t)
	}

	fm := template.FuncMap{
		"po": func(i int) int {
			return i + 1
		},
	}

	baseHTML = "templates/index.html"
	tmpl := template.Must(template.New("layout").Funcs(fm).ParseFiles(layout, baseHTML))

	err = tmpl.ExecuteTemplate(w, "layout", tasks)
	if err != nil {
		log.Println(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	baseHTML = "templates/create.html"
	tmpl, err := template.ParseFiles(baseHTML, layout)
	if err != nil {
		log.Fatalf("Template parsing failed: %v", err)
	}
	tmpl.Execute(w, nil)
}

func submitTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	urgent := r.FormValue("urgent")

	if title == "" || content == "" || urgent == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	t := NewTask(title, content)

	q := `INSERT INTO tasks (title, content, urgent)
	VALUES($1, $2, $3)`
	_, err := db.Exec(q, t.Title, t.Content, urgent)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	q := `DELETE FROM tasks WHERE id=$1`
	_, err := db.Exec(q, id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func singleTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	q := `SELECT * FROM tasks WHERE id=$1`
	row := db.QueryRow(q, id)

	t := Task{}
	err := row.Scan(&t.ID, &t.Title, &t.Content, &t.IsUrgent)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	baseHTML = "templates/edit.html"
	tmpl, err := template.ParseFiles(baseHTML, layout)
	if err != nil {
		log.Fatalf("Template parsing failed: %v", err)
	}
	tmpl.Execute(w, t)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	urgent := r.FormValue("urgent")
	if id == "" || title == "" || content == "" || urgent == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	q := `UPDATE tasks SET title=$1, content=$2, urgent=$3 WHERE id=$4`
	_, err := db.Exec(q, title, content, urgent, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
