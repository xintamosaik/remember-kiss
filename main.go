package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
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
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(globalTODOsInMemory)
}

func loadFromJSON() {
	data, err := os.ReadFile("data.json")
	if err != nil {
		return
	} // file might not exist yet
	json.Unmarshal(data, &globalTODOsInMemory)
}

func PageAdd(w http.ResponseWriter, r *http.Request) {
	page := struct {
		filename string
		title    string
		files    []string
	}{
		filename: "add.html",
		title:    "Add Item",
		files:    []string{"frame.html", "add.html"},
	}

	tmpl := template.Must(template.ParseFiles(page.files...))
	data := struct{ Title string }{Title: page.title}

	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", page.filename, err)
	}
}

func PageEdit(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	item, exists := globalTODOsInMemory[key]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title string
		Short string
		Long  string
	}{
		Title: "Edit Item",
		Short: item.Short,
		Long:  item.Long,
	}

	tmpl := template.Must(template.ParseFiles("frame.html", "edit.html"))
	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", "frame.html", err)
	}
}

func PageIndex(w http.ResponseWriter, r *http.Request) {
	page := struct {
		filename string
		title    string
		files    []string
	}{
		filename: "index.html",
		title:    "TODO",
		files:    []string{"frame.html", "main.html"},
	}

	tmpl := template.Must(template.ParseFiles(page.files...))

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
		Title: page.title,
		Items: items,
	}

	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", page.filename, err)
	}
}

func Todo(w http.ResponseWriter, r *http.Request) {
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
}

func addTodo(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Current TODOs: %+v\n", globalTODOsInMemory)

	// Handle form submission logic here
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Current TODOs: %+v\n", globalTODOsInMemory)

	// Handle form submission logic here
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("GET /index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.css")
	})

	http.HandleFunc("GET /", PageIndex)
	http.HandleFunc("GET /add.html", PageAdd)
	http.HandleFunc("GET /edit.html", PageEdit)

	// http.HandleFunc("GET /todo/", Todo)
	http.HandleFunc("POST /add", addTodo)
	http.HandleFunc("POST /update", updateTodo)

	log.Println("Serving on http://localhost" + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
