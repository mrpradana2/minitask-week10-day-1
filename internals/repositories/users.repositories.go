package repositories

import (
	"context"
	"tikcitz-app/internals/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Repository add user
func (u *UserRepository) UserRegister(ctx context.Context, email string, password string, role string) (pgconn.CommandTag, error) {
	// menambahkan user baru dengan mengembalikan id user baru
	queryUser := "INSERT INTO users (email, password, role) VALUES ($1, $2, &3) RETURNING id"
	var userID int
	err := u.db.QueryRow(ctx, queryUser, email, password, role).Scan(&userID)

	if err != nil {
		return pgconn.CommandTag{}, err
	}

	// default value
	first_name := ""
	last_name := ""
	phone_number := ""
	photo_path := ""
	title := ""
	point := 0
	
	// menambahkan baris baru untuk data profile user baru namun dengan default value
	queryProfile := "INSERT INTO profile (user_id, first_name, last_name, phone_number, photo_path, title, point) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	cmd, err := u.db.Exec(ctx, queryProfile, userID, first_name, last_name, phone_number, photo_path, title, point)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

// Repository user login
func (u *UserRepository) UserLogin(ctx context.Context, auth models.UsersStruct) (models.UsersStruct, error) {
	// mengambil data user dari DB berdasarkan email
	query := "SELECT email, password FROM users WHERE email = $1"
	var result models.UsersStruct
	err := u.db.QueryRow(ctx, query, auth.Email).Scan(&result.Email, &result.Password)
	if err != nil {
		return models.UsersStruct{}, err
	}

	return result, nil
}

// Repository get rpofile by id
func (u *UserRepository) GetProfileById(ctx context.Context, idInt int) (models.ProfileStruct, error) {
	query := "SELECT user_id, phone_number, first_name, last_name, photo_path, title, point FROM profile WHERE user_id = $1"
	values := []any{idInt}
	var result models.ProfileStruct
	if err := u.db.QueryRow(ctx, query, values...).Scan(&result.User_Id, &result.Phone_number, &result.First_name, &result.Last_name, &result.Photo_path, &result.Title, &result.Point); err != nil {
		return models.ProfileStruct{}, err
	}
	return result, nil
}

// Repository update profile
func (u *UserRepository) UpdateProfile(ctx context.Context, updateProfile models.ProfileStruct, idInt int) (pgconn.CommandTag, error) {
	query := "UPDATE profile SET first_name = $1, last_name = $2, phone_number = $3, photo_path = $4, title = $5, modified_at = $6 WHERE user_id = $7"
	values := []any{updateProfile.First_name, updateProfile.Last_name, updateProfile.Phone_number, updateProfile.Photo_path, updateProfile.Title, time.Now(), idInt}
	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}