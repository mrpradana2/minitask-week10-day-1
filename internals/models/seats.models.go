package models

import "time"

type CinemaStruct struct {
	Id          int    `json:"id"`
	Cinema_name string `json:"cinema_name"`
	Image_path  string `json:"image_path"`
}

type SeatsStruct struct {
	Id   int    `json:"id"`
	Seat string `json:"seat,omitempty"`
}

type ResultSeat struct {
	Id     int `json:"id"`
	Title  string
	Cinema string
	Date   time.Time
	Time time.Time
	Seats []string
}