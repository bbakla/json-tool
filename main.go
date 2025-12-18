package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"json_formatter/handlers"
)

// PageData holds values passed into the template.
type PageData struct {
	RawInput   string
	Key        string
	Value      string
	Formatted  string
	Matches    []string
	KeyMatches []string
	Error      string
}

func main() {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl := template.Must(template.ParseFiles(tmplPath))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, tmpl, PageData{})
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				renderTemplate(w, tmpl, PageData{Error: "failed to parse form"})
				return
			}
			rawJSON := r.FormValue("jsonInput")
			key := r.FormValue("key")
			value := r.FormValue("value")

			formatted, matches, keyMatches, err := handlers.Process(rawJSON, key, value)
			data := PageData{
				RawInput:   rawJSON,
				Key:        key,
				Value:      value,
				Formatted:  formatted,
				Matches:    matches,
				KeyMatches: keyMatches,
			}
			if err != nil {
				data.Error = err.Error()
			}
			renderTemplate(w, tmpl, data)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
