package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"

	vueglue "github.com/torenware/vite-go"
)

var environment string
var assets string
var jsEntryPoint string

var vueData *vueglue.VueGlue

//go:embed frontend
var dist embed.FS

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func pageWithAVue(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./test-template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, vueData)
}

func main() {

	flag.StringVar(&environment, "env", "development", "development|production")
	flag.StringVar(&assets, "assets", "frontend", "location of javascript files. dist for production.")
	flag.StringVar(&jsEntryPoint, "entryp", "src/main.js", "relative path of the entry point of the js app.")
	flag.Parse()

	// We pass the file system with the built Vue
	// program, and the path from the root of that
	// file system to the "assets" directory.

	var config vueglue.ViteConfig

	config.Environment = environment
	config.AssetsPath = assets
	config.EntryPoint = jsEntryPoint
	//config.FS = os.DirFS(assets)
	config.FS = dist

	if environment == "production" {
		config.URLPrefix = "/assets/"
	} else if environment == "development" {
		config.URLPrefix = "/src/"
	} else {
		log.Fatalln("illegal environment setting")
	}

	glue, err := vueglue.NewVueGlue(&config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	vueData = glue

	// Set up our router
	mux := http.NewServeMux()

	// Set up a file server for our assets.
	fsHandler, err := glue.FileServer()
	if err != nil {
		log.Println("could not set up static file server", err)
		return
	}
	mux.Handle("/src/", fsHandler)
	mux.HandleFunc("/", pageWithAVue)

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
