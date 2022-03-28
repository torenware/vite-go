package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	vueglue "github.com/torenware/vite-go"
)

//go:embed "dist"
var dist embed.FS

var vueData *vueglue.VueGlue

func pageWithAVue(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./test-template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, vueData)
}
func main() {

	// We pass the file system with the built Vue
	// program, and the path from the root of that
	// file system to the "assets" directory.
	glue, err := vueglue.NewVueGlue(dist, "dist")
	if err != nil {
		log.Fatalln(err)
		return
	}
	vueData = glue

	// Set up our router
	mux := http.NewServeMux()
	mux.HandleFunc("/", pageWithAVue)

	// Add a static file server for our Vue related elements
	sub, err := fs.Sub(dist, "dist/assets")
	if err != nil {
		log.Fatalln(err)
	}
	assetServer := http.FileServer(http.FS(sub))
	mux.Handle("/assets/", http.StripPrefix("/assets", assetServer))

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
