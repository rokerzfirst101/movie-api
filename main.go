package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"demoapi/models"
)

var movieList map[string]models.Movie

func main() {
	movieList = make(map[string]models.Movie)

	http.HandleFunc("/movies", CreateMovie)
	http.HandleFunc("/movies/", MovieHandler)
	http.ListenAndServe(":9000", nil)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var movie models.Movie

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "Can't Read Body", http.StatusBadRequest)
		return
	}

	parseErr := json.Unmarshal(body, &movie)
	if parseErr != nil {
		fmt.Print(parseErr)
		http.Error(w, "Invalid Body:", http.StatusBadRequest)
		return
	}
	_, isOk := movieList[movie.ID]
	if isOk {
		http.Error(w, "Movie Already Exists", http.StatusBadRequest)
		return
	}

	movieList[movie.ID] = movie
	val, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	io.WriteString(w, string(val[:]))
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := strings.Split(r.URL.String(), "/")[2]

	switch r.Method {
	case http.MethodGet:
		value, ok := movieList[id]
		if !ok {
			http.Error(w, "Movie does not exist.", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(value)
	case http.MethodDelete:
		_, ok := movieList[id]
		if !ok {
			fmt.Fprintf(w, "GET %q", html.EscapeString(r.URL.Path))
			return
		}
		delete(movieList, id)
	case http.MethodPut:
		var movie models.Movie
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Can't Read Body", http.StatusBadRequest)
			return
		}

		parseErr := json.Unmarshal(body, &movie)
		if parseErr != nil {
			fmt.Print(parseErr)
			http.Error(w, "Invalid Body:", http.StatusBadRequest)
			return
		}
		movieList[id] = movie
	}
}
