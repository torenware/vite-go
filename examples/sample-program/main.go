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
	"regexp"
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

func serveOneFile(w http.ResponseWriter, r *http.Request, uri, contentType string) {
	strippedURI := uri[1:]
	buf, err := fs.ReadFile(vueData.DistFS, strippedURI)
	if err != nil {
		// Try public dir
		buf, err = fs.ReadFile(vueData.DistFS, "public/"+strippedURI)
	}

	// If we ended up nil, render the file out.
	if err == nil {
		// not an error; letting the error case fall through
		w.Header().Add("Content-Type", contentType)
		w.Write(buf)
		return
	}

	// Otherwise, we cannot handle it, so 404 it is.
	w.WriteHeader(http.StatusNotFound)
}

func pageWithAVue(w http.ResponseWriter, r *http.Request) {
	// since we are using vite test pages, they assume
	// that the vite logo is at /vite.svg.  Let's handle
	// that case here. Since we're only serving one
	// page via Go, we'll pass svg, icongs and jpg
	// on to the dev server via a redirect
	re := regexp.MustCompile(`^/([^.]+)\.(svg|ico|jpg)$`)
	matches := re.FindStringSubmatch(r.RequestURI)
	if matches != nil {
		if vueData.Environment == "development" {
			log.Printf("vite logo requested")
			url := vueData.BaseURL + r.RequestURI
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		} else {
			// production; we need to render this ourselves.
			var contentType string
			switch matches[2] {
			case "svg":
				contentType = "image/svg+xml"
			case "ico":
				contentType = "image/x-icon"
			case "jpg":
				contentType = "image/jpeg"
			}

			serveOneFile(w, r, r.RequestURI, contentType)
			return
		}

	}

	// our go page, which will host our javascript.
	t, err := template.ParseFiles("./test-template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, vueData)
}

func main() {
	var config vueglue.ViteConfig

	flag.StringVar(&config.Environment, "env", "development", "development|production")
	flag.StringVar(&config.JSProjectPath, "assets", "", "location of javascript files.")
	flag.StringVar(&config.AssetsPath, "dist", "", "dist directory relative to the JS project directory.")
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
		config.FS = os.DirFS("frontend")
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
