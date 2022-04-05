package vueglue

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
)

// type ViteConfig passes info needed to generate the library's
// output.
type ViteConfig struct {

	// FS is the filesystem where assets can be loaded.
	FS fs.FS

	// Environment (development|production)
	Environment string

	//AssetsPath (typically dist/assets for prod, src for dev)
	AssetsPath string

	// URLPrefix (assets/ for prod, src/ for dev)
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

	// I use a 'data-entryp' attribute to find what
	// components to load. Lookup is in the entry point JS.
	// This makes the info easily available in templates.
	MountPoint string

	// An embed that points to the Vue/Vite dist
	// directory.
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
	var err error
	var glue *VueGlue
	glue = &VueGlue{}

	glue.Environment = config.Environment
	glue.DistFS = config.FS

	if config.Environment == "production" {
		// Get the manifest file
		manifestFile := filepath.Join(config.AssetsPath, "manifest.json")
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

	output, _ := json.MarshalIndent(glue, "", "  ")
	fmt.Println(string(output))

	tags, err := glue.RenderTags()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(tags)

	return glue, nil

}
