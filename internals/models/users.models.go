package models

import "mime/multipart"

type ProfileStruct struct {
	User_Id      int    `db:"user_id" json:"user_id" form:"user_id,omitempty"`
	First_name   string `db:"first_name" json:"first_name" form:"first_name,omitempty"`
	Last_name    string `db:"last_name" json:"last_name" form:"last_name,omitempty"`
	Phone_number string `db:"phone_number" json:"phone_number" form:"phone_number,omitempty"`
	Photo_path   *multipart.FileHeader `db:"photo_path" json:"-" form:"photo_path,omitempty"`
	PhotoPath    string `db:"photo_profile" json:"photo_profile" form:"photo_profile,omitempty"`
	Title        string `db:"title" json:"title" form:"title,omitempty"`
	Point        int    `db:"point" json:"point" form:"point,omitempty"`
	NewPassword     string `json:"-" form:"new_password" binding:"min=8"`
	ComfirmPassword string `json:"-" form:"confirm_password" binding:"min=8"`
}

type UsersStruct struct {
	Id       int    `db:"id,omitempty" json:"id,omitempty"`
	Email    string `db:"email" json:"email" form:"email" binding:"required,email"`
	Password string `db:"password" json:"password" form:"password" binding:"min=8"`
	Role     string `db:"role,omitempty" json:"role,omitempty"`
}

type SignupPayload struct {
	UsersStruct
	ProfileStruct
}