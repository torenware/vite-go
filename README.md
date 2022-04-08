# Vite Integration For Go

A simple module that lets you serve your Vue 3 project from a Go-based web server.  You build your project, tell Go where to find the `dist/` directory, and the module figures out how to load the generated Vue application into a web page.

## Installation

```shell
go get github.com/torenware/vite-go

```

To upgrade to the current version:

```shell
go get -u github.com/torenware/vite-go@latest 
```


## Getting It Into Your Go Project

You need to expose the `dist/` directory so Go can find your generated assets for the Vue project, and the `manifest.json` file that describes it. You may need to change your `vite.config.ts` file to make sure the manifest file is generated, and to put the `dist` directory where Go needs it to be. Here's what I'm using:

```typescript
/**
 * @type {import('vite').UserConfig}
 */
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'cmd/web/dist',
    sourcemap: true,
    manifest: true,
    rollupOptions: {
      input: {
        main: 'src/main.ts',
      },
    },
  },
});
```  

This, however, more than you need. A minimal config file would be:

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    manifest: "manifest.json",
  },
})

```

The essential piece here is the vue plugin and the `build.manifest` line, since `vite-go` needs the manifest file to be present in order to work correctly.


Here's some pseudo sample code that uses the go 1.16+ embedding feature for the production build, and a regular disk directory (`frontend` in our case) as a development directory:

```golang

package main

import (
  "embed"
  "html/template"
  "net/http"

  vueglue "github.com/torenware/vite-go"
)

//go:embed "dist"
var dist embed.FS

var vueGlue *vueglue.VueGlue

func main() {
    
    // This:
    
    // Production configuration.
	config := &vueglue.ViteConfig{
		Environment: "production",
		AssetsPath:  "dist",
		EntryPoint:  "src/main.js",
		FS:          dist,
	}

    // OR this:
    
    // Development configuration
	config := &vueglue.ViteConfig{
		Environment: "development",
		AssetsPath:  "frontend",
		EntryPoint:  "src/main.js",
		FS:          os.DirFS("frontend"),
	}

  // Parse the manifest and get a struct that describes
  // where the assets are.
  glue, err := vueglue.NewVueGlue(config)
  if err != nil {
    //bail!
  }
  vueGlue = glue
  
  // and set up your routes and start your server....
  
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
  // Now you can pass the glue object to an HTML template
  ts, err := template.ParseFiles("path/to/your-template.tmpl")
  if err != nil {
  	// better handle this...
  }
  ts.Execute(respWriter, vueGlue)

}


```

You will also need to serve your javascript, css and images used by your javascript code to the web. You can use a solution like [`http.FileServer`](https://pkg.go.dev/net/http#FileServer), or the wrapper the library implements that configures this for you:

```golang
   // using the standard library's multiplexer:
	mux := http.NewServeMux()

	// Set up a file server for our assets.
	fsHandler, err := glue.FileServer()
	if err != nil {
		log.Println("could not set up static file server", err)
		return
	}
	mux.Handle("/src/", fsHandler)

```

Some router implementations may alternatively require you to do something more like:

```golang
// chi router
mux := chi.NewMux()

...

mux.Handle("/src/*", fsHandler)

```

YMMV :-)

## Templates

Your template gets the needed tags and links something like this:


```HTML
<!doctype html>
<html lang='en'>
{{ $vue := . }}
    <head>
        <meta charset='utf-8'>
        <title>Home - Vue Loader Test</title>
        
        {{ if $vue }}
          {{ $vue.RenderTags }}
        {{ end }}
        
    </head>
    <body>
      <div id="app"></div>
      
    </body>
  </html>
      
 
```

The sample program in `examples/sample-program` has much more detail, and actually runs.

## Caveats

This code is a proof of concept, and while it works in my sample application, it may not work for you :-) I've posted the code so people can see it, and kick the tires on it. I think you'll find it useful.



Copyright © 2022 Rob Thorne

[MIT License](https://github.com/torenware/go-tooling-for-vue/blob/8999977a5bffb8f0630740220c576b550a7115e9/LICENSE.md)
<hr>
