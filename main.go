package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	tmpl := template.Must(template.ParseFiles(
		"frame.html",
		"content.html",
	))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":   "Hello world",
			"Heading": "Welcome",
			"Body":    "This content comes from layout.html",
		}

		if err := tmpl.ExecuteTemplate(w, "frame", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Serving on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
