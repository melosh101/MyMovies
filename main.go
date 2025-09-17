//go:build dev
// +build dev

package main

import (
	"encoding/json"
	"html/template"
	"htmxBackend/templates"
	"htmxBackend/tmdb"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	mux.HandleFunc("/api/movies/list", handleMovieList)

	mux.HandleFunc("/api/handle-post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World! post"))
	})

	mux.HandleFunc("/hx/{action}", handleHXRequest)

	log.Println("Listening on :8080")
	handler := mainHandler(mux)
	err = http.ListenAndServe(":8080", handler)
	// err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleHXRequest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Hx-Request") != "true" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Expected HX action"))
	}

	action := r.PathValue("action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Expected action"))
	}
	log.Println("action:", action)
	switch action {
	case "nextPage":
		{
			pageStr := r.URL.Query().Get("page")
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			movies, err := tmdb.ListMovies(tmdb.ListMoviesParms{Page: page})
			if err != nil {
				log.Println(err)
				return
			}
			data := templates.MoviePageData{
				Films: movies.Results,
			}
			templ := template.Must(template.ParseFiles("public/templates/index.gohtml"))
			err = templ.ExecuteTemplate(w, "movies", data)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("HX Action Not Found"))
}

func handleMovieList(w http.ResponseWriter, r *http.Request) {

	url, _ := url2.Parse(r.URL.String())
	params, _ := url2.ParseQuery(url.RawQuery)

	page, _ := strconv.Atoi(params.Get("page"))

	if page == 0 {
		page = 1
	}

	listParams := tmdb.ListMoviesParms{
		Page: page,
	}

	res, err := tmdb.ListMovies(listParams)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
		return
	}

}
