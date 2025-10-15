package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const PORT = ":8080"

// This generates 3 static html files index.html, add.html and edit.html in the folder public
func regenerateHTML() {
	tmpl := template.Must(template.ParseFiles("frame.html"))

	pages := []struct {
		filename    string
		title       string
		contentFile string
	}{
		{"public/index.html", "TODO", "main.html"},
		{"public/add.html", "Add Item", "add.html"},
		{"public/edit.html", "Edit Item", "edit.html"},
	}

	for _, page := range pages {
		contentBytes, err := os.ReadFile(page.contentFile)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", page.contentFile, err)
		}
		content := struct {
			Title   string
			Content template.HTML
		}{
			Title:   page.title,
			Content: template.HTML(string(contentBytes)),
		}
		fs, err := os.Create(page.filename)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", page.filename, err)
		}
		err = tmpl.ExecuteTemplate(fs, "frame", content)
		fs.Close()
		if err != nil {
			log.Fatalf("Failed to render %s: %v", page.filename, err)
		}
	}
}
func main() {

	http.Handle("/", http.FileServer(http.Dir("public")))
	regenerateHTML() // will be in a different process later but for development it's ok here

	log.Println("Serving on http://localhost" + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
