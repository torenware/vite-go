package vueglue

import (
	"os"
	"strings"
	"testing"
)

func TestDevServer(t *testing.T) {

	config := &ViteConfig{
		Environment: "development",
		AssetsPath:  "tests/testdata",
		URLPrefix:   "/",
		FS:          os.DirFS("testdata"),
		EntryPoint:  "main.ts",
		ViteVersion: "3",
	}

	// check default tag generated
	glue, err := initializeVueGlue(config)
	if err != nil {
		t.Fatalf("lib did not initialize: %s", err)
	}
	shouldContain := "main.ts"
	if glue.MainModule != shouldContain {
		t.Fatalf("dev tag looks wrong. expected %s, got %s", shouldContain, glue.MainModule)
	}

	tags, err := glue.RenderTags()
	if err != nil {
		t.Fatalf("tags did not render: %s", err)
	}
	// t.Logf("tags: %s", string(tags))
	shouldContain = "http://localhost:5173/main.ts"
	if !strings.Contains(string(tags), shouldContain) {
		t.Fatalf("tags did not contain '%s'", shouldContain)
	}

	// change defaults
	config.DevServerPort = "3001"
	glue, err = initializeVueGlue(config)
	if err != nil {
		t.Fatalf("could not parse config: %s", err)
	}

	tags, err = glue.RenderTags()
	if err != nil {
		t.Fatalf("non-default tags not rendered: %s", err)
	}
	shouldContain = glue.BaseURL + "/main.ts"
	if !strings.Contains(string(tags), shouldContain) {
		t.Fatalf("tags did not contain '%s'", shouldContain)
	}

}
