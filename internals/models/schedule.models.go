package models

import "time"

type ScheduleStruct struct {
	Id   int `json:"id"`
	Title string `json:"title"`
	Cinema string `json:"cinema"`
	Date time.Time	`json:"-"`
	Time time.Time `json:"-"`
	DateStr string `json:"date"`
	TimeStr string `json:"time"`
	Price int `json:"price"`
	Location string `json:"location"`
}

