package models

import "time"

type MoviesStruct struct {
	Id				int	`json:"id,omitempty"`
	Title           string    `json:"title" form:"title"`
	Image_path      string    `json:"image_path" form:"image_path"`
	Overview        string    `json:"overview" form:"overview"`
	Release_date    time.Time `json:"release_date" form:"release_date"`
	Director_name   string    `json:"director_name" form:"director_name"`
	Duration        int       `json:"duration" form:"duration"`
	Casts           []string  `json:"casts" form:"casts"`
	Genres          []string  `json:"genres" form:"genres"`
	Status_movie_id int `json:"status_movie_id" form:"status_movie_id"` //hapus nanti
	Status_movie string `json:"status_movie" form:"status_movie"` //hapus nanti
	ScheduleS
}

type ScheduleS struct {
	Cinema_ids []int `json:"cinema_ids" form:"cinema_ids"`
	Movie_id int `json:"movie_id" form:"movie_id"`
	Location string	`json:"location" form:"location"`
	Date time.Time `json:"date" form:"date"`
	Times []time.Time `json:"time" form:"time"`
	Price int `json:"price" form:"price"`
}