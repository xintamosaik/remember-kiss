package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"encoding/json"
)

type TodoItem struct {
	ID        string
	Short     string
	Long      string
	CreatedAt string
	UpdatedAt string
}

type TodoList struct {
	Items []TodoItem
}

const PORT = ":3000"

var globalTODOsInMemory = make(map[string]TodoItem)
func saveToJSON() error {
    f, err := os.Create("data.json")
    if err != nil { return err }
    defer f.Close()

    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    return enc.Encode(globalTODOsInMemory)
}

func loadFromJSON() {
    data, err := os.ReadFile("data.json")
    if err != nil { return } // file might not exist yet
    json.Unmarshal(data, &globalTODOsInMemory)
}

func regenerateIndex() {
	tmpl := template.Must(template.ParseFiles("frame.html", "main.html"))

	loadFromJSON()
	// Convert map to slice for template rendering
	items := make([]TodoItem, 0, len(globalTODOsInMemory))
	for _, item := range globalTODOsInMemory {
		items = append(items, item)
	}

	data := struct {
		Title string
		Items []TodoItem
	}{
		Title: "TODO",
		Items: items,
	}

	f, err := os.Create("public/index.html")
	if err != nil {
		log.Fatalf("Failed to create index.html: %v", err)
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, "frame", data)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}
}

func regenerateHTML() {
	pages := []struct {
		filename string
		title    string
		files    []string
	}{
		{"public/add.html", "Add Item", []string{"frame.html", "add.html"}},
		{"public/edit.html", "Edit Item", []string{"frame.html", "edit.html"}},
	}

	for _, page := range pages {
		tmpl := template.Must(template.ParseFiles(page.files...))
		data := struct{ Title string }{Title: page.title}

		f, err := os.Create(page.filename)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", page.filename, err)
		}
		defer f.Close()

		err = tmpl.ExecuteTemplate(f, "frame", data)
		if err != nil {
			log.Fatalf("Failed to render %s: %v", page.filename, err)
		}
	}
}
func main() {

	http.Handle("/", http.FileServer(http.Dir("public")))
	regenerateHTML() // will be in a different process later but for development it's ok here
	regenerateIndex()
	// API:POST:/add
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		short := r.FormValue("short")
		log.Println(short)

		// create a linux timestamp for the key
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		item := TodoItem{
			ID:        timestamp,
			Short:     short,
			Long:      "",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		globalTODOsInMemory[timestamp] = item
		saveToJSON()
		regenerateIndex()
		log.Printf("Current TODOs: %+v\n", globalTODOsInMemory)

		// Handle form submission logic here
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	log.Println("Serving on http://localhost" + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
