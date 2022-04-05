package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	vueglue "github.com/torenware/vite-go"
)

var environment string
var assets string
var jsEntryPoint string

var vueData *vueglue.VueGlue

func GuardedFileServer(prefix, stripPrefix, serveDir string) http.Handler {

	handler := func(w http.ResponseWriter, r *http.Request) {
		prefixLen := len(stripPrefix)
		rest := r.URL.Path[prefixLen:]
		parts := strings.Split(rest, "/")
		// We want to prevent dot files
		if parts[len(parts)-1][:1] == "." {
			//force a relative link.
			log.Printf("Found dotfile or dir %s", parts[0])
			http.NotFound(w, r)
			return
		}
		fileServer := http.StripPrefix(stripPrefix, http.FileServer(http.Dir(serveDir)))
		fileServer.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}

func ServeVueAssets(mux *http.ServeMux, prefix, stripPrefix, serveDir string) error {
	assetServer := GuardedFileServer(prefix, stripPrefix, serveDir)
	mux.Handle(prefix, logRequest(assetServer))
	return nil
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log called")
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
	config.FS = os.DirFS(assets)

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
	mux.HandleFunc("/", pageWithAVue)

	err = ServeVueAssets(mux, config.URLPrefix, "/", config.AssetsPath)
	if err != nil {
		log.Println("setting up FS failed:", err)
	}

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
