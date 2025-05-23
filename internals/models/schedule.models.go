package models

import "time"

type ScheduleStruct struct {
	Id   int `json:"id"`
	MovieId   int `json:"movie_id"`
	Title string `json:"title"`
	Cinema string `json:"cinema"`
	CinemaPathImage string `json:"cinema_img"`
	Date time.Time	`json:"-"`
	Time time.Time `json:"time"`
	// Time2 time.Time `json:"waktu"`
	DateStr string `json:"date"`
	Location string `json:"location"`
	Price int `json:"price"`
}

