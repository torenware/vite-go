package vueglue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
)

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Type            string            `json:"type"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencis"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (vg *ViteConfig) getViteVersion() (string, error) {
	// If it's set, use it.
	if vg.ViteVersion != "" {
		return vg.ViteVersion, nil
	}

	// If not set, try and find package.json
	buf, err := fs.ReadFile(vg.FS, "package.json")
	if err != nil {
		return "", err
	}

	content := PackageJSON{}
	err = json.Unmarshal(buf, &content)
	if err != nil {
		return "", err
	}

	vite, ok := content.DevDependencies["vite"]
	if !ok {
		return "", errors.New("no devdep/vite")
	}

	re, err := regexp.Compile(`^\^(\d+)`)
	if err != nil {
		return "", err
	}

	match := re.FindStringSubmatch(vite)
	if match == nil {
		return "", errors.New("no match")
	}

	vg.ViteVersion = match[1]
	return match[1], nil
}

func (vg *ViteConfig) setDevelopmentDefaults() {
	version, err := vg.getViteVersion()
	if err != nil {
		vg.ViteVersion = DEFAULT_VITE_VERSION
		version = vg.ViteVersion
	}

	if vg.DevServerPort == "" {
		if version == "2" {
			vg.DevServerPort = DEFAULT_PORT_V2
		} else {
			vg.DevServerPort = DEFAULT_PORT_V3
		}
	}

	if vg.DevServerDomain == "" {
		vg.DevServerDomain = "localhost"
	}

}

func (vg *ViteConfig) buildDevServerBaseURL() string {
	protocol := "http"
	if vg.HTTPS {
		protocol = "https"
	}

	return fmt.Sprintf(
		"%s://%s:%s",
		protocol,
		vg.DevServerDomain,
		vg.DevServerPort,
	)

}
