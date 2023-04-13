package vueglue

import (
	"embed"
	"errors"
	"io/fs"
)

// constants
const (
	DEFAULT_VITE_VERSION = "3"
	DEFAULT_PORT_V2      = "3000"
	DEFAULT_PORT_V3      = "5173"
)

// type ViteConfig passes info needed to generate the library's
// output.
type ViteConfig struct {

	// FS is the filesystem where assets can be loaded.
	FS fs.FS

	// DevDefaults is best guess for defaults
	DevDefaults *JSAppParams `json:"-"`

	// Environment (development|production). In development mode,
	// the package sets up hot reloading. In production, the
	// package builds the Vue/Vuex production files and embeds them
	// in the Go app.
	Environment string

	// JSProjectPath is where your JS project is relative to the
	// root of your project. Default: frontend
	JSProjectPath string

	// JSInExternalDir denotes that you keep your JS project source
	// in a folder located external to your go project source.
	// Default: false
	JSInExternalDir bool

	//AssetsPath relative to the JSProjectPath. Empty for dev, dist for prod
	AssetsPath string

	// "2" or "3". If not set, we try to guess by looking
	// at package.json
	ViteVersion string

	// DevServerDomain is what domain the dev server appears on.
	// Default is localhost.
	DevServerDomain string

	// DevServerPort is what port the dev server will appear on.
	// Default depends upon the ViteVersion.
	DevServerPort string

	// HTTPS is whether the dev server is encrypted or not.
	// Default is false.
	HTTPS bool

	// URLPrefix (/assets/ for prod, /src/ for dev)
	URLPrefix string

	// DevServer is the URL to use for the Vite dev server.
	// Default is "http://localhost:3000".
	// DevServer string

	// Platform (vue|react|svelte) is the target platform.
	// Default is "vue"
	Platform string

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

	// BaseURL is the base URL for the dev server.
	// Default is http://localhost:5173
	BaseURL string

	// JS Dependencies / Vendor libs
	Imports []string

	// Bundled CSS
	CSSModule []string

	// Target JS Platform
	Platform string

	// A file system or embed that points to the Vue/Vite dist
	// directory (production) or the javascript src directory
	// (development)
	DistFS fs.FS

	// DevServer is the URI of the Vite development server
	DevServer string

	// JSProjectPath is the location of the JS project.
	JSProjectPath string

	// AssetPath is the relative path from the JSDirectory.
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

	correctedFS, err := correctEmbedFS(config.FS, config.JSProjectPath)
	if err != nil {
		return nil, err
	}

	if config.Environment == "production" {
		err := config.SetProductionDefaults()
		if err != nil {
			return nil, err
		}

		// Get the manifest file
		manifestFile := config.AssetsPath + "/manifest.json"
		contents, err := fs.ReadFile(correctedFS, manifestFile)
		if err != nil {
			return nil, err
		}
		glue, err = ParseManifest(contents)
		if err != nil {
			return nil, err
		}

	} else {
		err := config.SetDevelopmentDefaults()
		if err != nil {
			return nil, err
		}
		glue.BaseURL = config.buildDevServerBaseURL()
		glue.MainModule = config.EntryPoint
	}

	glue.Environment = config.Environment
	glue.JSProjectPath = config.JSProjectPath
	glue.AssetPath = config.AssetsPath
	glue.Platform = config.Platform
	glue.DistFS = correctedFS

	return glue, nil
}
