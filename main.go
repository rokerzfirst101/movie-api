package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"demoapi/models"

	"github.com/gorilla/mux"
)

var movieList map[int]models.Movie

func main() {
	movieList = make(map[int]models.Movie)

	r := mux.NewRouter()

	r.HandleFunc("/movies", CreateMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{id}", GetById).Methods(http.MethodGet)
	r.HandleFunc("/movies/{id}", UpdateById).Methods(http.MethodPut)
	r.HandleFunc("/movies/{id}", DeleteById).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9000",
	}
	log.Fatalln(srv.ListenAndServe())
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var movie models.Movie

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parseErr := json.Unmarshal(body, &movie)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, isOk := movieList[movie.ID]
	if isOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movieList[movie.ID] = movie
	val, err := json.Marshal(movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(val)
}

func GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, ok := movieList[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(value)
}

func UpdateById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, exists := movieList[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var movie models.Movie

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parseErr := json.Unmarshal(body, &movie)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	movieList[id] = movie
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

func DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := movieList[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	delete(movieList, id)
}
