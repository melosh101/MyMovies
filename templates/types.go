package templates

import "htmxBackend/tmdb"

type MoviePageData struct {
	Films      []tmdb.MovieListResult
	NowPlaying map[string][]tmdb.MovieListResult
}
