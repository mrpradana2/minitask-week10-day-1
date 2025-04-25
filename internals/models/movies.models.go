package models

import "time"

type MoviesStruct struct {
	Title           string    `json:"title" form:"title"`
	Image_path      string    `json:"image_path" form:"image_path"`
	Overview        string    `json:"overview" form:"overview"`
	Release_date    time.Time `json:"release_date" form:"release_date"`
	Director_name   string    `json:"director_name" form:"director_name"`
	Duration        int       `json:"duration" form:"duration"`
	Casts           []string  `json:"casts" form:"casts"`
	Status_movie_id int       `json:"status_movie_id" form:"status_movie_id"`
	Genres          []string  `json:"genres" form:"genres"`
	Status_movie 	string	`json:"status_movie" form:"status_movie"`
}