package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"sort"
)

type TodoItem struct {
	ID        string `json:"ID"`
	Short     string `json:"Short"`
	Long      string `json:"Long"`
	Done      bool   `json:"Done"`
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
	}
	json.Unmarshal(data, &globalTODOsInMemory)
}

func PageAdd(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frame.html", "add.html"))
	data := struct{ Title string }{Title: "Add Item"}

	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", "add.html", err)
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
		Key   string
		Short string
		Long  string
	}{
		Title: "Edit Item",
		Key:   key,
		Short: item.Short,
		Long:  item.Long,
	}

	tmpl := template.Must(template.ParseFiles("frame.html", "edit.html"))
	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", "frame.html", err)
	}
}

func PageDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	item, exists := globalTODOsInMemory[key]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title string
		Key   string
		Short string
		Long  string
	}{
		Title: "Delete Item",
		Key:   key,
		Short: item.Short,
		Long:  item.Long,
	}

	tmpl := template.Must(template.ParseFiles("frame.html", "delete.html"))
	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", "frame.html", err)
	}
}

func PageIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frame.html", "main.html"))

	items := make([]TodoItem, 0, len(globalTODOsInMemory))
	for _, item := range globalTODOsInMemory {
		items = append(items, item)
	}

	// Sort items by CreatedAt timestamp (optional)
	sort.Slice(items, func(i, j int) bool {
	 	return items[i].CreatedAt < items[j].CreatedAt
	})

	data := struct {
		Title string
		Items []TodoItem
	}{
		Title: "TODO",
		Items: items,
	}

	err := tmpl.ExecuteTemplate(w, "frame", data)
	if err != nil {
		log.Fatalf("Failed to render %s: %v", "index.html", err)
	}
}

func Todo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todo/"):]

	item, exists := globalTODOsInMemory[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Deleting item %s\n", key)

	_, exists := globalTODOsInMemory[key]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(globalTODOsInMemory, key)
	saveToJSON()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	
	log.Printf("Toggling item %s\n", key)
	item, exists := globalTODOsInMemory[key]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item.Done = !item.Done
	item.UpdatedAt = time.Now().Format(time.RFC3339)

	globalTODOsInMemory[key] = item
	saveToJSON()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	loadFromJSON()
	http.HandleFunc("GET /index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.css")
	})

	http.HandleFunc("GET /", PageIndex)
	
	http.HandleFunc("GET /add.html", PageAdd)
	http.HandleFunc("POST /add", addTodo)

	http.HandleFunc("GET /edit.html", PageEdit)
	http.HandleFunc("POST /update", updateTodo)

	http.HandleFunc("GET /delete.html", PageDelete)
	http.HandleFunc("POST /delete", deleteTodo)

	http.HandleFunc("GET /toggle/", toggleTodo)

	log.Println("Serving on http://localhost" + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
