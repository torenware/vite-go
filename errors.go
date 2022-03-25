package vueglue

import (
	"errors"
)

var (
	ErrNoEntryPoint        = errors.New("manifest lacked entry point")
	ErrNoInputFile         = errors.New("expected import file name")
	ErrManifestBadlyFormed = errors.New("manifest has unexpected format")
	ErrManifestDNF         = errors.New("vue distribution directory not found")
)
