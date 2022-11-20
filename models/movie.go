package models

type Movie struct {
	ID       int     `json: "id"`
	Name     string  `json: "name"`
	Genre    string  `json: "genre"`
	Rating   float64 `json: "rating"`
	Plot     string  `json: "plot"`
	Released bool    `json: "released"`
}
