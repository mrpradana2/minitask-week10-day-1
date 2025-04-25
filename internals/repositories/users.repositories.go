package repositories

import (
	"context"
	"tikcitz-app/internals/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) UserRegister(ctx *gin.Context, newDataUser models.SignupPayload) (pgconn.CommandTag, error) {

	queryUser := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var userID int
	err := u.db.QueryRow(ctx.Request.Context(), queryUser, newDataUser.Email, newDataUser.Password).Scan(&userID)

	if err != nil {
		return pgconn.CommandTag{}, err
	}

	queryProfile := "INSERT INTO profile (user_id, modified_at) VALUES ($1, $2)"
	cmd, err := u.db.Exec(ctx.Request.Context(), queryProfile, userID, time.Now())
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

func (u *UserRepository) UserLogin(ctx *gin.Context, auth models.UsersStruct) (models.UsersStruct, error) {
	// mengambil data user dari DB
	query := "SELECT email, password FROM users WHERE email = $1"
	var result models.UsersStruct
	err := u.db.QueryRow(ctx.Request.Context(), query, auth.Email).Scan(&result.Email, &result.Password)
	if err != nil {
		return models.UsersStruct{}, err
	}

	return result, nil
}

func (u *UserRepository) GetProfileById(ctx *gin.Context, idInt int) (models.ProfileStruct, error) {
	query := "SELECT phone_number, first_name, last_name, photo_path, title FROM profile WHERE user_id = $1"
	values := []any{idInt}
	var result models.ProfileStruct
	if err := u.db.QueryRow(context.Background(), query, values...).Scan(&result.Phone_number, &result.First_name, &result.Last_name, &result.Photo_path, &result.Title); err != nil {
		return models.ProfileStruct{}, err
	}
	return result, nil
}

func (u *UserRepository) UpdateProfile(ctx *gin.Context, updateProfile models.ProfileStruct, idInt int) (pgconn.CommandTag, error) {
	query := "UPDATE profile SET first_name = $1, last_name = $2, phone_number = $3, photo_path = $4, title = $5, modified_at = $6 WHERE user_id = $7"
	values := []any{updateProfile.First_name, updateProfile.Last_name, updateProfile.Phone_number, updateProfile.Photo_path, updateProfile.Title, time.Now(), idInt}
	cmd, err := u.db.Exec(ctx.Request.Context(), query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}