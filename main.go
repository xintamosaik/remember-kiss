package main

import (
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("<p>Hello, world!</p>"))
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("public"))) // kinda only if dev mode but for now, whatever

	http.HandleFunc("/api/hello", helloHandler);

	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
