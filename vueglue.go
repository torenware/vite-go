package vueglue

import (
	"embed"
	"io/fs"
	"log"
)

// type ViteConfig passes info needed to generate the library's
// output.
type ViteConfig struct {

	// FS is the filesystem where assets can be loaded.
	FS fs.FS

	// Environment (development|production). In development mode,
	// the package sets up hot reloading. In production, the
	// package builds the Vue/Vuex production files and embeds them
	// in the Go app.
	Environment string

	//AssetsPath (typically dist for prod, and your Vue project
	// directory for dev)
	AssetsPath string

	// URLPrefix (/assets/ for prod, /src/ for dev)
	URLPrefix string

	// Entry point: as configured in vite.config.js. Typically
	// src/main.js or src/main.ts.
	EntryPoint string
}

// type VueGlue summarizes a manifest file, and points to the assets.
type VueGlue struct {

	// Environment. This controls whether the library will
	// configure the host for hot updating, or whether it
	// needs to configure loading of a dist/ directory.
	Environment string

	// Entry point for JS
	MainModule string

	// JS Dependencies / Vendor libs
	Imports []string

	// Bundled CSS
	CSSModule []string

	// A file system or embed that points to the Vue/Vite dist
	// directory (production) or the javascript src directory
	// (development)
	DistFS fs.FS
}

// ParseManifest imports and parses a manifest returning a glue object.
func ParseManifest(contents []byte) (*VueGlue, error) {
	var testRslt manifestTarget
	glue, err := testRslt.parseWithoutReflection(contents)
	if err != nil {
		return nil, err
	}
	return glue, nil
}

// NewVueGlue finds the manifest in the supplied file system
// and returns a glue object.
func NewVueGlue(config *ViteConfig) (*VueGlue, error) {
	var glue *VueGlue
	glue = &VueGlue{}

	glue.Environment = config.Environment
	glue.DistFS = config.FS

	if config.Environment == "production" {
		// embed behaves a little strange: it does
		// not set the top level dir as the "current"
		// dir for the FS. So give it a clue.
		// @see https://github.com/golang/go/issues/43431
		prefix := ""
		if _, ok := config.FS.(embed.FS); ok {
			log.Println("we are using an embed")
			prefix = config.AssetsPath + "/"
		}
		// Get the manifest file
		manifestFile := prefix + "manifest.json"
		contents, err := fs.ReadFile(config.FS, manifestFile)
		if err != nil {
			return nil, err
		}
		glue, err = ParseManifest(contents)
		if err != nil {
			return nil, err
		}
	} else {
		// all we need for hot updating.
		glue.MainModule = config.EntryPoint
	}

	return glue, nil
}
