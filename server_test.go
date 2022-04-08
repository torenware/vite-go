package vueglue

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//go:embed testdata
var embedTest embed.FS

func InitializeVueGlue(config *ViteConfig) (*VueGlue, error) {
	glue, err := NewVueGlue(config)
	return glue, err
}

func StartTestServer(glue *VueGlue) (*httptest.Server, error) {
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

func BootStrapServer(config *ViteConfig) (*httptest.Server, error) {
	if config == nil {
		config = &ViteConfig{
			Environment: "development",
			AssetsPath:  "tests/testdata",
			URLPrefix:   "/",
			FS:          os.DirFS("testdata"),
			EntryPoint:  "server.js",
		}
	}

	glue, err := InitializeVueGlue(config)
	if err != nil {
		return nil, err
	}

	server, err := StartTestServer(glue)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func TestInitLib(t *testing.T) {
	config := &ViteConfig{
		Environment: "development",
		AssetsPath:  "tests/testdata",
		URLPrefix:   "/",
		FS:          os.DirFS("testdata"),
		EntryPoint:  "server.js",
	}

	glue, err := InitializeVueGlue(config)
	if err != nil {
		t.Fatalf("Library failed to initialize: %s", err)
	}

	if glue == nil {
		t.Fatalf("No glue was returned")
	}

	if glue.MainModule != "server.js" {
		t.Fatalf("Expected main module to be %s, got %s", "server.ts", glue.MainModule)
	}

}

func TestServerHandler(t *testing.T) {
	config := &ViteConfig{
		Environment: "development",
		AssetsPath:  "testdata",
		URLPrefix:   "/",
		FS:          os.DirFS("testdata"),
		EntryPoint:  "server.js",
	}
	glue, err := InitializeVueGlue(config)
	if err != nil {
		t.Fatalf("no glue! %s", err)
	}
	_, err = glue.FileServer()
	if err != nil {
		t.Fatalf("no handler was returned: %s", err)
	}

	srv, err := StartTestServer(glue)
	if err != nil {
		t.Fatalf("server did not bootstrap: %s", err)
	}
	defer srv.Close()

	url := srv.URL

	response, err := http.Head(url)
	if err != nil {
		t.Fatalf("could not ping server: %s", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("HEAD / got %d, expected %d", response.StatusCode, 200)
	}

}

func TestFileVisibility(t *testing.T) {
	srv, err := BootStrapServer(nil)
	if err != nil {
		t.Fatalf("could not bootstrap test server: %s", err)
	}
	defer srv.Close()

	var dataList = []struct {
		Path   string
		Status int
	}{
		{"", 200},
		{"index.html", 200},
		{"regfile.txt", 200},
		{"not-there", 404},
		{".secret", 404},
		{".secdir/file.txt", 404},
		{"subdir/regfile.txt", 200},
		{"subdir/.env-file", 404},
	}

	base := srv.URL
	for _, item := range dataList {
		url := fmt.Sprintf("%s/%s", base, item.Path)
		response, err := http.Head(url)
		if err != nil {
			t.Errorf("%s: Error on Head %s", item.Path, err)
		} else {
			if response.StatusCode != item.Status {
				t.Errorf("%s: expected %d but got %d", item.Path, item.Status, response.StatusCode)
			}
		}
	}

}

func TestEmbedAccess(t *testing.T) {
	config := &ViteConfig{
		Environment: "development",
		AssetsPath:  "testdata",
		URLPrefix:   "/",
		FS:          embedTest,
		EntryPoint:  "server.js",
	}
	srv, err := BootStrapServer(config)
	if err != nil {
		t.Fatalf("could not bootstrap test server: %s", err)
	}
	defer srv.Close()

	var dataList = []struct {
		Path   string
		Status int
	}{
		{"", 200},
		{"index.html", 200},
		{"regfile.txt", 200},
		{"not-there", 404},
		{".secret", 404},
		{".secdir/file.txt", 404},
		{"subdir/regfile.txt", 200},
		{"subdir/.env-file", 404},
	}

	base := srv.URL
	for _, item := range dataList {
		url := fmt.Sprintf("%s/%s", base, item.Path)
		response, err := http.Head(url)
		if err != nil {
			t.Errorf("%s: Error on Head %s", item.Path, err)
		} else {
			if response.StatusCode != item.Status {
				t.Errorf("%s: expected %d but got %d", item.Path, item.Status, response.StatusCode)
			}
		}
	}

}
