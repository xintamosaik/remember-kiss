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
	ID        string `json:"ID"`
	Short     string `json:"Short"`
	Long      string `json:"Long"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

type TodoList struct {
	Items []TodoItem `json:"Items"`
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
	// Load existing TODOs from JSON file
	regenerateHTML() // will be in a different process later but for development it's ok here
	regenerateIndex()

	http.Handle("/", http.FileServer(http.Dir("public")))


	// API:GET:/todo/{id} - for fetching a single item with fetch API
	http.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/todo/"):]
		log.Printf("Fetching item with ID: %s\n", id)
		item, exists := globalTODOsInMemory[id]
		if !exists {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		log.Printf("Found item: %+v\n", item)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)	
	})

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
	// API:POST:/update
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		key := r.FormValue("key")
		short := r.FormValue("short")
		long := r.FormValue("long")
		log.Printf("Updating item %s: %s / %s\n", key, short, long)

		item, exists := globalTODOsInMemory[key]
		if !exists {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		item.Short = short
		item.Long = long
		item.UpdatedAt = time.Now().Format(time.RFC3339)
		globalTODOsInMemory[key] = item
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
