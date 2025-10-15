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
	content := struct {
		Title   string
		Content template.HTML
	}{
		Title:   "TODO",
		Content: template.HTML("<main>Hi</main>"),
	}
	// write to file
	fs, err := os.Create("public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	err = tmpl.ExecuteTemplate(fs, "frame", content)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {


	http.Handle("/", http.FileServer(http.Dir("public")))
	regenerateHTML()

	log.Println("Serving on http://localhost" + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
