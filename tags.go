package vueglue

import (
	"bytes"
	"html/template"
)

// RenderTags genarates the HTML tags that link a rendered
// Go template with any Vue assets that need to be loaded.
func (vg *VueGlue) RenderTags() (template.HTML, error) {
	var tags string

	if vg.Environment == "development" {
		if vg.Platform == "react" {
			// react requires some extra help to load
			tags += `
    <script src="/src/preamble.js"></script>
            `
		}
		tags += `
    <script type="module" src="{{.DevServer}}/{{ .MainModule }}"></script>
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
	tmpl.Execute(&buffer, vg)

	return template.HTML(buffer.String()), nil
}
