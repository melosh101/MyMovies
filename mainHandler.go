//go:build !dev

package main

import (
	"embed"
	"html/template"
	"htmxBackend/templates"
	"htmxBackend/tmdb"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed public
var public embed.FS

func mainHandler(apiMux *http.ServeMux) http.HandlerFunc {
	subFS, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatalf("Error getting sub-filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(subFS))

	log.Println("using prod static handler")

	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" || r.URL.Path == "" {
			log.Println("serving /")
			res, err := tmdb.GetPopular()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			nowPlaying, err := tmdb.GetNowPlaying()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("%v\n", nowPlaying)
			templ := template.Must(template.ParseFS(subFS, "index.gohtml"))

			d := templates.MoviePageData{
				Films: res.Results,
				NowPlaying: map[string][]tmdb.MovieListResult{
					"Films": nowPlaying.Results[:8],
				},
			}

			err = templ.Execute(w, d)
			if err != nil {
				log.Println("error during render:", err)
			}
			return
		}

		filePath := strings.TrimPrefix(r.URL.Path, "/")

		// Check if the file exists in the embedded filesystem.
		// If the file does not exist, an error will be returned.
		file, err := fs.Stat(subFS, filePath)
		if err != nil {
			apiMux.ServeHTTP(w, r)
			return
		}
		// Check for the existence of the file or an `index.gohtml` at the root.
		if err == nil && file != nil && !file.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		log.Println("file not found")
		log.Println(filePath)

		// This handles root requests and any other URL path that doesn't resolve to a file
		// (e.g., /about, /dashboard), which is a common pattern for SPAs.
		http.ServeFileFS(w, r, subFS, "index.gohtml")
		log.Println("using fallback static handler")
	}
}
