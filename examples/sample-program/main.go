package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	vueglue "github.com/torenware/vite-go"
)

// this is not for vite, but to help our
// makefile stop the process:
var pidFile string
var pidDeleteChan chan os.Signal

func waitForSignal() {
	if pidFile == "" {
		return
	}
	pidDeleteChan = make(chan os.Signal, 1)
	signal.Notify(pidDeleteChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-pidDeleteChan
		fmt.Println("Deleted pid file")
		_ = os.Remove(pidFile)
		os.Exit(0)

	}()
}

var vueData *vueglue.VueGlue

//go:embed frontend
var dist embed.FS

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func serveOneFile(uri string, w http.ResponseWriter, r *http.Request) {
	// serving this out of our embed
	path := "frontend"
	contentType := "image/svg+xml"
	if uri == "/vite.svg" {
		path += "/public/vite.svg"
		buf, err := fs.ReadFile(dist, path)
		if err == nil {
			// not an error; letting the error case fall through
			w.Header().Add("Content-Type", contentType)
			w.Write(buf)
			return
		}
	}

	// Otherwise, we cannot handle it, so 404 it is.
	w.WriteHeader(http.StatusNotFound)
}

func pageWithAVue(w http.ResponseWriter, r *http.Request) {
	// since we are using vite test pages, they assume
	// that the vite logo is at /vite.svg.  Let's handle
	// that case here.
	if r.RequestURI == "/vite.svg" {
		log.Printf("vite logo requested")
		serveOneFile(r.RequestURI, w, r)
		return
	}

	t, err := template.ParseFiles("./test-template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, vueData)
}

func main() {
	var config vueglue.ViteConfig

	flag.StringVar(&config.Environment, "env", "development", "development|production")
	flag.StringVar(&config.AssetsPath, "assets", "", "location of javascript files. dist for production.")
	flag.StringVar(&config.EntryPoint, "entryp", "", "relative path of the entry point of the js app.")
	flag.StringVar(&config.Platform, "platform", "", "vue|react|svelte")
	flag.StringVar(&pidFile, "pid", "", "location of optional pid file.")
	flag.Parse()

	// save away our pid if we need to use a makefile to stop
	// this process. You don't need this to use Vite, but it does
	// make this test program easier to use.
	if pidFile != "" {
		pid := strconv.Itoa(os.Getpid())
		_ = ioutil.WriteFile(pidFile, []byte(pid), 0644)
		waitForSignal()
	}
	defer func() {
		if pidFile != "" {
			err := os.Remove(pidFile)
			if err != nil {
				log.Printf("could not delete pid file: %v", err)
			}
		}
	}()

	// We pass the file system with the built Vue
	// program, and the path from the root of that
	// file system to the "assets" directory.
	if config.EntryPoint == "production" {
		config.FS = os.DirFS("frontend/dist")
	} else {
		// Use the embed.
		config.FS = dist
	}
	//

	if config.Environment == "production" {
		config.URLPrefix = "/assets/"
	} else if config.Environment == "development" {
		log.Printf("pulling defaults using package.json")
	} else {
		log.Fatalln("unsupported environment setting")
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
	mux.Handle(config.URLPrefix, fsHandler)
	mux.Handle("/", logRequest(http.HandlerFunc(pageWithAVue)))

	log.Println("Starting server on :4000")
	generatedConfig, _ := json.MarshalIndent(config, "", "  ")
	log.Println("Generated Configuration:\n", string(generatedConfig))
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
