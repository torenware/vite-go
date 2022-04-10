package vueglue

import (
	"embed"
	"errors"
	"io/fs"
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

	// DevServer is the URL to use for the Vite dev server.
	// Default is "http://localhost:3000".
	DevServer string

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

	// DevServer is the URI of the Vite development server
	DevServer string

	// AssetPath is the path from the root of DistFS. This
	// allows us to check if the FS is correctly "pointed"
	// to the diretory containing the assets.
	AssetPath string

	// Debug mode
	Debug bool
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

// If we have an embedded FS, modify it to point to the
// requested assets directory
func correctEmbedFS(embedded fs.FS, assetsPath string) (fs.FS, error) {

	// embed behaves a little strange: it does
	// not set the top level dir as the "current"
	// dir for the FS. This is almost never what you
	// want, so we correct for this
	//
	// @see https://github.com/golang/go/issues/43431
	//
	if _, ok := embedded.(embed.FS); ok {
		// Make sure someone has not already taken a sub:
		_, err := fs.ReadDir(embedded, assetsPath)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if err == nil {
			// uncorrected FS, so take its subdir
			embedded, err = fs.Sub(embedded, assetsPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return embedded, nil
}

// NewVueGlue finds the manifest in the supplied file system
// and returns a glue object.
func NewVueGlue(config *ViteConfig) (*VueGlue, error) {
	var glue *VueGlue
	glue = &VueGlue{}

	correctedFS, err := correctEmbedFS(config.FS, config.AssetsPath)
	if err != nil {
		return nil, err
	}

	if config.Environment == "production" {
		// Get the manifest file
		manifestFile := "manifest.json"
		contents, err := fs.ReadFile(correctedFS, manifestFile)
		if err != nil {
			return nil, err
		}
		glue, err = ParseManifest(contents)
		if err != nil {
			return nil, err
		}

	} else {
		glue.MainModule = config.EntryPoint
		if config.DevServer == "" {
			glue.DevServer = "http://localhost:3000"
		} else {
			glue.DevServer = config.DevServer
		}
	}

	glue.Environment = config.Environment
	glue.AssetPath = config.AssetsPath
	glue.DistFS = correctedFS

	return glue, nil
}
