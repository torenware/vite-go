package vueglue

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
)

// type VueGlue summarizes a manifest file, and points to the assets.
type VueGlue struct {

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
	DistFS fs.ReadFileFS
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
func NewVueGlue(dist fs.ReadFileFS, pathToDist string) (*VueGlue, error) {

	if !fs.ValidPath(pathToDist) {
		return nil, ErrManifestDNF
	}

	// Get the manifest file
	manifestFile := filepath.Join(pathToDist, "manifest.json")
	contents, err := dist.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}
	glue, err := ParseManifest(contents)
	if err != nil {
		return nil, err
	}
	glue.DistFS = dist

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
