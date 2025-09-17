//go:build dev

package main

import (
	"htmxBackend/templates"
	"htmxBackend/tmdb"
	"log"
	"net/http"
	"text/template"
)

// mainHandler serves files from the local filesystem during development.
func mainHandler(apiMux *http.ServeMux) http.Handler {
	fs := http.FileServer(http.Dir("./public"))
	log.Println("using templates")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// First, check if the request is for an API endpoint.
		if len(r.URL.Path) > 4 && r.URL.Path[:4] == "/api" {
			apiMux.ServeHTTP(w, r)
			return
		}
		if len(r.URL.Path) > 3 && r.URL.Path[:3] == "/hx" {
			apiMux.ServeHTTP(w, r)
			return
		}

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
			templ := template.Must(template.ParseFiles("public/index.gohtml"))

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

		log.Println("serving /", r.URL.Path)
		// Otherwise, serve from the filesystem.
		fs.ServeHTTP(w, r)
	})
}

var baseImgUrl = "https://image.tmdb.org/t/p/"

func GetPosterUrl(posterPath string) string {
	return baseImgUrl + "w500" + posterPath
}
