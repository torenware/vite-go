package vueglue

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
)

func initializeVueGlue(config *ViteConfig) (*VueGlue, error) {
	if config == nil {
		config = &ViteConfig{
			Environment:   "development",
			JSProjectPath: "testdata",
			URLPrefix:     "/",
			FS:            os.DirFS("testdata"),
			EntryPoint:    "main.js",
			ViteVersion:   "2",
		}
	}
	glue, err := NewVueGlue(config)
	return glue, err
}

func startTestServer(glue *VueGlue) (*httptest.Server, error) {
	handler, err := glue.FileServer()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	server := httptest.NewServer(mux)
	if server == nil {
		return nil, errors.New("did not get server instance")
	}

	return server, nil
}

func bootStrapServer(config *ViteConfig) (*httptest.Server, error) {

	glue, err := initializeVueGlue(config)
	if err != nil {
		return nil, err
	}

	server, err := startTestServer(glue)
	if err != nil {
		return nil, err
	}
	return server, nil
}
