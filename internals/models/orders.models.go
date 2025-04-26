package models

type OrdersStruct struct {
	User_id            int    `json:"user_id" form:"user_id"`
	Movie_id           int    `json:"movie_id" form:"movie_id"`
	Total_price        int    `json:"total_price" form:"total_price"`
	Full_name          string `json:"full_name" form:"full_name"`
	Email              string `json:"email" form:"email"`
	Phone_number       string `json:"phone_number" form:"phone_number"`
	Payment_methode_id int    `json:"payment_methode_id" form:"payment_methode_id"`
	Paid               bool   `json:"paid" form:"paid"`
	Date               string `json:"date" form:"date"`
	Time               string `json:"time" form:"time"`
	Cinema_id          int    `json:"cinema_id" form:"cinema_id"`
}