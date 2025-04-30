package models

type ProfileStruct struct {
	User_Id      int    `db:"user_id" json:"user_id" form:"user_id,omitempty"`
	First_name   string `db:"first_name" json:"first_name" form:"first_name,omitempty"`
	Last_name    string `db:"last_name" json:"last_name" form:"last_name,omitempty"`
	Phone_number string `db:"phone_number" json:"phone_number" form:"phone_number,omitempty"`
	Photo_path   string `db:"photo_path" json:"photo_path" form:"photo_path,omitempty"`
	Title        string `db:"title" json:"title" form:"title,omitempty"`
	Point        int    `db:"point" json:"point" form:"point,omitempty"`
}

type UsersStruct struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type SignupPayload struct {
	UsersStruct
	ProfileStruct
}