package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var serverPort = ":" + os.Getenv("PORT")

var tmpl *template.Template

func init() {
	fm := template.FuncMap{
		"po": func(i int) int {
			return i + 1
		},
	}

	tmpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gohtml"))
}

func main() {
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

	if err := tmpl.ExecuteTemplate(w, "index.gohtml", tasks); err != nil {
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "create.gohtml", nil); err != nil {
		log.Fatal(err)
	}
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

	if err := tmpl.ExecuteTemplate(w, "edit.gohtml", nil); err != nil {
		log.Fatal(err)
	}
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
