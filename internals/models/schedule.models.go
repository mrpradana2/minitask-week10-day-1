package models

import "time"

type ScheduleStruct struct {
	Id   int `json:"movie_id"`
	Title string `json:"title"`
	Cinema string `json:"cinema"`
	CinemaPathImage string `json:"cinema_img"`
	Date time.Time	`json:"date"`
	Time []time.Time `json:"time"`
	DateStr string `json:"datestr"`
	Location string `json:"location"`
	Price int `json:"price"`
}

