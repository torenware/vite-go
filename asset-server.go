package vueglue

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// FileServer is a customized version of http.FileServer
// that can handle either an embed.FS or a os.DirFS fs.FS.
// Since development directories used for hot updates
// can contain dot files (potentially with sensitive
// information) the code checks to make sure that dot files
// are not served.
func (vg *VueGlue) FileServer() (http.Handler, error) {
	// First, make sure if our fs.FS is from an embed.FS,
	// that we adjust where the FS is "pointing".
	target, err := correctEmbedFS(vg.DistFS, vg.AssetPath)
	if err != nil {
		return nil, err
	}

	// Prevent directory listings
	wrapped := wrapperFS{
		FS: target,
	}

	handler := vg.guardedFileServer(wrapped)

	return handler, nil
}

// guardedFileServer wraps http.FileServer with filtering
// code that checks for dot files and other potentially
// sensitive material a static file server should not
// render over.
//
// We assume that the fs.Dir's top level is pointed at the contents
// of where the assets are, and not its parent directory as would
// typically be the case for an embed.FS instance.
func (vg *VueGlue) guardedFileServer(serveDir fs.FS) http.Handler {
	stripPrefix := "/"
	handler := func(w http.ResponseWriter, r *http.Request) {
		prefixLen := len(stripPrefix)
		rest := r.URL.Path[prefixLen:]
		parts := strings.Split(rest, "/")

		// Now walk the parts and make sure none of them are
		// either "hidden" files or directories.
		for _, stem := range parts {
			if len(stem) > 0 && stem[:1] == "." {
				http.NotFound(w, r)
				return
			}
		}

		if vg.Debug {
			log.Println("entered FS", r.URL.Path)
			dir, err := fs.ReadDir(serveDir, ".")
			if err != nil {
				log.Println("could not read the asset dir", err)
				http.NotFound(w, r)
				return
			}

			for _, item := range dir {
				log.Println(item.Name())
			}

		}

		fileServer := http.StripPrefix(stripPrefix, http.FileServer(http.FS(serveDir)))
		fileServer.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}

// Wrapper file system to prevent listing of directories
// @see https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings

type wrapperFS struct {
	FS fs.FS
}

// Open implements the fs.FS interface for wrapperFS
func (wrpr wrapperFS) Open(path string) (fs.File, error) {
	f, err := wrpr.FS.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		// Have an index file or go home!
		index := filepath.Join(path, "index.html")
		if _, err := wrpr.FS.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
