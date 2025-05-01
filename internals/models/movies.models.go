package models

import "time"

// m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name)

type MoviesStruct struct {
	Id				int	`json:"id,omitempty"`
	Title           string    `json:"title" form:"title"`
	Release_date    time.Time `json:"release_date" form:"release_date"`
	Overview        string    `json:"overview" form:"overview"`
	Image_path      string    `json:"image_path" form:"image_path"`
	Duration        int       `json:"duration" form:"duration"`
	Director_name   string    `json:"director_name" form:"director_name"`
	Genres          []string  `json:"genres" form:"genres"`
	Casts           []string  `json:"casts" form:"casts"`
	Status_movie_id int `json:"status_movie_id,omitempty" form:"status_movie_id"` //hapus nanti
	Status_movie string `json:"status_movie,omitempty" form:"status_movie"` //hapus nanti
	ScheduleS
}

type ScheduleS struct {
	Cinema_ids []int `json:"cinema_ids,omitempty" form:"cinema_ids"`
	Movie_id int `json:"movie_id,omitempty" form:"movie_id"`
	Location string	`json:"location,omitempty" form:"location"`
	Date time.Time `json:"date,omitempty" form:"date,omitempty"`
	Times []time.Time `json:"time,omitempty" form:"time"`
	Price int `json:"price,omitempty" form:"price"`
}