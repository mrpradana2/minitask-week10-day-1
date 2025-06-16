package models

import (
	"mime/multipart"
	"time"
)

// m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name)

type MoviesStruct struct {
	Id				int	        `json:"id,omitempty"`
	Title           string      `json:"title" form:"title" binding:"required"`
	Release_date    time.Time   `json:"release_date" form:"release_date" binding:"required"`
	Overview        string      `json:"overview" form:"overview" binding:"required"`
	Duration        int         `json:"duration" form:"duration" binding:"required"`
	Director_name   string      `json:"director_name" form:"director_name" binding:"required"`
	Genres          []string    `json:"genres" form:"genres" binding:"required"`
	Casts           []string    `json:"casts" form:"casts" binding:"required"`
	Image_movie     string      `json:"image_movie"`
	TotalSales      int         `json:"total_sales,omitempty"`
	Image_path      *multipart.FileHeader `json:"image_path,omitempty" form:"image_path,omitempty" binding:"required"` 
	Cinema_ids      []int       `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
	Movie_id        int         `json:"movie_id,omitempty" form:"movie_id,omitempty"`
	Location        string 	    `json:"location,omitempty" form:"location,omitempty"`
	Date            time.Time   `json:"-" form:"date,omitempty"`
	Times           []time.Time `json:"time,omitempty" form:"time,omitempty"`
	Price           int         `json:"price,omitempty" form:"price,omitempty"`
}

type UpdateMoviesStruct struct {
	Id				int	        `json:"id,omitempty"`
	Title           string      `json:"title" form:"title" binding:"required"`
	Release_date    time.Time   `json:"release_date" form:"release_date" binding:"required"`
	Overview        string      `json:"overview" form:"overview" binding:"required"`
	Duration        int         `json:"duration" form:"duration" binding:"required"`
	Director_name   string      `json:"director_name" form:"director_name" binding:"required"`
	Genres          []string    `json:"genres" form:"genres" binding:"required"`
	Casts           []string    `json:"casts" form:"casts" binding:"required"`
	Image_movie     string      `json:"image_movie"`
	TotalSales      int         `json:"total_sales,omitempty"`
	Image_path      *multipart.FileHeader `json:"image_path,omitempty" form:"image_path,omitempty"` 
	Old_Image_path  string      `json:"old_image_path,omitempty" form:"old_image_path,omitempty"` 
	Cinema_ids      []int       `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
	Movie_id        int         `json:"movie_id,omitempty" form:"movie_id,omitempty"`
	Location        string 	    `json:"location,omitempty" form:"location,omitempty"`
	Date            time.Time   `json:"-" form:"date,omitempty"`
	Times           []time.Time `json:"time,omitempty" form:"time,omitempty"`
	Price           int         `json:"price,omitempty" form:"price,omitempty"`
}

type RequestMoviesStr struct {
	Title           string      `json:"title" form:"title" binding:"required"`
	Release_date    time.Time   `json:"release_date" form:"release_date" binding:"required"`
	Overview        string      `json:"overview" form:"overview" binding:"required"`
	Duration        int         `json:"duration" form:"duration" binding:"required"`
	Director_name   string      `json:"director_name" form:"director_name" binding:"required"`
	Genres          []string    `json:"genres" form:"genres" binding:"required"`
	Casts           []string    `json:"casts" form:"casts" binding:"required"`
	Cinema_ids      []int       `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
	Location        string 	    `json:"location,omitempty" form:"location,omitempty"`
	Date            time.Time   `json:"-" form:"date,omitempty"`
	Times           []time.Time `json:"time,omitempty" form:"time,omitempty"`
	Price           int         `json:"price,omitempty" form:"price,omitempty"`
	Image_path      string      `json:"image_path,omitempty" form:"image_path,omitempty" binding:"required"` 
}

type RequestUpdateMoviesStr struct {
	Title           string      `json:"title" form:"title" binding:"required"`
	Release_date    time.Time   `json:"release_date" form:"release_date" binding:"required"`
	Overview        string      `json:"overview" form:"overview" binding:"required"`
	Duration        int         `json:"duration" form:"duration" binding:"required"`
	Director_name   string      `json:"director_name" form:"director_name" binding:"required"`
	Genres          []string    `json:"genres" form:"genres" binding:"required"`
	Casts           []string    `json:"casts" form:"casts" binding:"required"`
	Cinema_ids      []int       `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
	Location        string 	    `json:"location,omitempty" form:"location,omitempty"`
	Date            time.Time   `json:"-" form:"date,omitempty"`
	Times           []time.Time `json:"time,omitempty" form:"time,omitempty"`
	Price           int         `json:"price,omitempty" form:"price,omitempty"`
	Image_path      string      `json:"image_path,omitempty" form:"image_path,omitempty" binding:"required"`
	Old_Image_path  string      `json:"old_image_path,omitempty" form:"old_image_path,omitempty"` 
}