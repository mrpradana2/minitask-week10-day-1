package models

import (
	"mime/multipart"
	"time"
)

// m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name)

type MoviesStruct struct {
	Id				int	        `json:"id,omitempty"`
	Title           string      `json:"title" form:"title"`
	Release_date    time.Time   `json:"release_date" form:"release_date"`
	Overview        string      `json:"overview" form:"overview"`
	Duration        int         `json:"duration" form:"duration"`
	Director_name   string      `json:"director_name" form:"director_name"`
	Genres          []string    `json:"genres" form:"genres"`
	Casts           []string    `json:"casts" form:"casts"`
	Image_movie     string      `json:"image_movie"`
	TotalSales      int         `json:"total_sales"`
	Image_path      *multipart.FileHeader `json:"image_path,omitempty" form:"image_path,omitempty"` 
	Cinema_ids      []int       `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
	Movie_id        int         `json:"movie_id,omitempty" form:"movie_id,omitempty"`
	Location        string 	    `json:"location,omitempty" form:"location,omitempty"`
	Date            time.Time   `json:"-" form:"date,omitempty"`
	Times           []time.Time `json:"time,omitempty" form:"time,omitempty"`
	Price           int         `json:"price,omitempty" form:"price,omitempty"`
}