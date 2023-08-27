package vueglue

import (
	"bytes"
	"html/template"
)

// RenderTags genarates the HTML tags that link a rendered
// Go template with any Vue assets that need to be loaded.
func (vg *VueGlue) RenderTags() (template.HTML, error) {
	var tags string

	tplPayload := struct {
		*VueGlue
		ReactRefreshURL template.JSStr
	}{
		VueGlue:         vg,
		ReactRefreshURL: template.JSStr(vg.BaseURL + "/@react-refresh"),
	}

	if vg.Environment == "development" {
		if vg.Platform == "react" {
			// react requires some extra help to load
			tags += `<script type="module">
				import RefreshRuntime from {{.ReactRefreshURL}};
				RefreshRuntime.injectIntoGlobalHook(window)
				window.$RefreshReg$ = () => {}
				window.$RefreshSig$ = () => (type) => type
				window.__vite_plugin_react_preamble_installed__ = true
			</script>
            `
		}
		tags += `
    <script type="module" src="{{.BaseURL}}/{{ .MainModule }}"></script>
        `

	} else {
		tags += `
	<script type="module" crossorigin src="/{{ .MainModule }}"></script>
	{{ range .Imports }}
	<link rel="modulepreload" href="/{{.}}">
	{{ end }}
	{{ range .CSSModule }}
	<link rel="stylesheet" href="/{{.}}">
	{{ end }}
	`
	}
	tmpl, err := template.New("tags").Parse(tags)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	tmpl.Execute(&buffer, tplPayload)

	return template.HTML(buffer.String()), nil
}
