package models

import "mime/multipart"

type ProfileStruct struct {
	User_Id      int    `db:"user_id" json:"user_id" form:"user_id,omitempty"`
	First_name   string `db:"first_name" json:"first_name" form:"first_name,omitempty"`
	Last_name    string `db:"last_name" json:"last_name" form:"last_name,omitempty"`
	Phone_number string `db:"phone_number" json:"phone_number" form:"phone_number,omitempty"`
	PhotoPath    string `db:"photo_profile" json:"photo_profile" form:"photo_profile,omitempty"`
	Title        string `db:"title" json:"title" form:"title,omitempty"`
	Point        int    `db:"point" json:"point" form:"point,omitempty"`
	Email		 string `db:"email" json:"email,omitempty" form:"email,omitempty"`
	NewPassword     string `json:"new_password,omitempty" form:"new_password"`
	ConfirmPassword string `json:"confirm_password,omitempty" form:"confirm_password"`
}

type RequestProfileStruct struct {
	First_name   string `db:"first_name" json:"first_name" form:"first_name,omitempty"`
	Last_name    string `db:"last_name" json:"last_name" form:"last_name,omitempty"`
	Phone_number string `db:"phone_number" json:"phone_number" form:"phone_number,omitempty"`
	Title        string `db:"title" json:"title" form:"title,omitempty"`
	Email		 string `db:"email" json:"email,omitempty" form:"email,omitempty"`
	NewPassword     string `json:"new_password,omitempty" form:"new_password"`
	ConfirmPassword string `json:"confirm_password,omitempty" form:"confirm_password"`
}

type PhotoProfileStruct struct {
	Photo_path   *multipart.FileHeader `db:"photo" json:"-" form:"photo,omitempty"`
}

// swagger:model
type RequestPhotoProfileStruct struct {
	Photo_path   string `db:"photo" json:"photo" form:"photo,omitempty"`
}

type UsersStruct struct {
	Id       int    `db:"id,omitempty" json:"id,omitempty"`
	Email    string `db:"email" json:"email" form:"email" binding:"required,email"`
	Password string `db:"password" json:"password" form:"password"`
	Role     string `db:"role,omitempty" json:"role,omitempty"`
}

type UserLogin struct {
	Email    string `db:"email" json:"email" form:"email" binding:"required,email"`
	Password string `db:"password" json:"password" form:"password" binding:"required,min=8"`
}