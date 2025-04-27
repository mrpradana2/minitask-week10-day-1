package models

import (
	"time"
)

type OrdersStruct struct {
	Id *int `json:"id,omitempty" form:"id,omitempty"`
	User_id            int    `json:"user_id" form:"user_id"`
	Movie_id           int    `json:"movie_id,omitempty" form:"movie_id"`
	Total_price        int    `json:"total_price" form:"total_price"`
	Full_name          string `json:"full_name,omitempty" form:"full_name"`
	Email              string `json:"email,omitempty" form:"email"`
	Phone_number       string `json:"phone_number,omitempty" form:"phone_number"`
	Payment_methode_id int    `json:"payment_methode_id,omitempty" form:"payment_methode_id"`
	Paid               bool   `json:"paid" form:"paid"`
	Date               time.Time `json:"date" form:"date"`
	Time               string `json:"time" form:"time"`
	Cinema_id          int    `json:"cinema_id,omitempty" form:"cinema_id"`
	Cinema_path     string `json:"cinema_path"`
	Title           string `json:"title"`
	Payment_methode string `json:"payment_methode"`
	SeatId []int `json:"seat_id" form:"seat_id"`
}