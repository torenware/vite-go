package vueglue

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func TestDefaultParams(t *testing.T) {

	tstList := []struct {
		file       string
		platform   string
		viteVer    string
		devPort    string
		typescript bool
		entryPt    string
	}{
		{"package-preact-ts.json", "preact", "3", "5173", true, "src/main.tsx"},
		{"package-react-ts.json", "react", "3", "5173", true, "src/main.tsx"},
		{"package-react.json", "react", "3", "5173", false, "src/main.jsx"},
		{"package-svelte-ts.json", "svelte", "3", "5173", true, "src/main.ts"},
		{"package-svelte.json", "svelte", "3", "5173", false, "src/main.js"},
		{"package-vanilla-ts.json", "vanilla", "3", "5173", true, "src/main.ts"},
		{"package-vanilla.json", "vanilla", "3", "5173", false, "main.js"},
		{"package-vue-ts.json", "vue", "3", "5173", true, "src/main.ts"},
		{"package-vue-v2.json", "vue", "2", "3000", false, "src/main.js"},
		{"package-vue3.json", "vue", "3", "5173", false, "src/main.js"},
	}

	for _, test := range tstList {
		file := fmt.Sprintf("testdata/pkg-json/%s", test.file)
		buf, err := fs.ReadFile(os.DirFS("."), file)
		if err != nil {
			t.Errorf("%s could not be read", test.file)
		}

		pkg := PackageJSON{}
		err = json.Unmarshal(buf, &pkg)
		if err != nil {
			t.Errorf("%s could not be parsed", test.file)
		}

		params := analyzePackageJSON(&pkg)
		if params == nil {
			t.Errorf("%s could not be analyzed", test.file)
		}

		if params.ViteMajorVer != test.viteVer {
			t.Errorf(
				"%s: Vite Version expected %s, got %s",
				test.file,
				test.viteVer,
				params.ViteMajorVer,
			)
		}

		if params.HasTypeScript != test.typescript {
			t.Errorf(
				"%s: typescript expected %t, got %t",
				test.file,
				test.typescript,
				params.HasTypeScript,
			)
		}

		if params.PackageType != test.platform {
			t.Errorf(
				"%s: platform expected %s, got %s",
				test.file,
				test.platform,
				params.PackageType,
			)
		}

		if params.EntryPoint != test.entryPt {
			t.Errorf(
				"%s: Entry Point expected %s, got %s",
				test.file,
				test.entryPt,
				params.EntryPoint,
			)
		}

	}

}
