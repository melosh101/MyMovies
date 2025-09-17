package tmdb

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type AuthTransport struct {
	Transport http.RoundTripper
	Token     string
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqClone := req.Clone(req.Context())
	reqClone.Header.Set("Authorization", "Bearer "+t.Token)
	reqClone.Header.Set("Content-Type", "application/json")
	return t.Transport.RoundTrip(reqClone)
}

var token = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJjMzAyNjQxNThhODQ4M2Q5MTU1OTk5NzRhMzYzMzgxNyIsIm5iZiI6MTc1NzkyMjY0OS44NzIsInN1YiI6IjY4YzdjNTU5ZTE3NGI1OTYyOTcwZDcxMSIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.VxZc8Y2y-zGL9av_mpInLAQf4AWHD4GlmWE5tvWEbdk"
var BASE_URL = "https://api.themoviedb.org/3/"

func getClient() *http.Client {
	return &http.Client{
		Transport: &AuthTransport{
			Transport: http.DefaultTransport,
			Token:     token,
		},
	}
}

type ListMoviesParms struct {
	Page int `default:"1"`
}

func ListMovies(params ListMoviesParms) (MovieList, error) {

	if params.Page < 1 {
		return MovieList{}, errors.New("Page must be at least 1 or above")
	}

	client := getClient()
	reqParams := url.Values{}
	url := BASE_URL + "discover/movie"

	reqParams.Set("page", strconv.Itoa(params.Page))

	resp, err := client.Get(url + "?" + reqParams.Encode())

	if err != nil {
		log.Println("Error getting movies:", err)
		return MovieList{}, errors.New("failed to fetch movies, check server logs")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return MovieList{}, errors.New("failed to fetch movies during body read, check server logs")
	}
	var result MovieList
	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Println("Error unmarshalling movies:", err)
		return MovieList{}, errors.New("failed to fetch movies during unmarshalling, check server logs")
	}
	return result, nil
}

func GetNowPlaying() (NowPlayingList, error) {

	client := getClient()
	url := BASE_URL + "movie/now_playing"

	resp, err := client.Get(url)
	if err != nil {
		log.Println("Error getting nowplaying:", err)
		return NowPlayingList{}, errors.New("failed to fetch nowplaying, check server logs")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error getting nowplaying:", err)
		return NowPlayingList{}, errors.New("failed to read body in nowplaying, check server logs")
	}

	var result NowPlayingList

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error unmarshalling nowplaying:", err)
	}
	return result, nil

}

func GetPopular() (MovieList, error) {
	client := getClient()
	url := BASE_URL + "movie/popular"

	resp, err := client.Get(url)
	if err != nil {
		log.Println("Error getting nowplaying:", err)
		return MovieList{}, errors.New("failed to fetch nowplaying, check server logs")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error getting nowplaying:", err)
		return MovieList{}, errors.New("failed to read body in nowplaying, check server logs")
	}

	var result MovieList

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error unmarshalling nowplaying:", err)
	}
	return result, nil

}
