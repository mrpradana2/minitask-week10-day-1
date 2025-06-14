package repositories

import (
	"context"
	"log"
	"strings"
	"tikcitz-app/internals/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct{
	db *pgxpool.Pool
	rdb *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, rdb *redis.Client) *UserRepository {
	return &UserRepository{db: db, rdb: rdb}
}

// Repository add user
func (u *UserRepository) UserRegister(ctx context.Context, email string, password string, role string) (pgconn.CommandTag, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	defer tx.Rollback(ctx)

	// menambahkan user baru dengan mengembalikan id user baru
	queryUser := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	var userID int
	if err := tx.QueryRow(ctx, queryUser, email, password, role).Scan(&userID); err != nil {
		log.Println("[ERROR] : ", err.Error())
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
	queryProfile := "INSERT INTO profiles (user_id, first_name, last_name, phone_number, photo_path, title, point) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	cmd, err := tx.Exec(ctx, queryProfile, userID, first_name, last_name, phone_number, photo_path, title, point)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

// Repository user login
func (u *UserRepository) UserLogin(ctx context.Context, auth models.UsersStruct) (models.UsersStruct, models.ProfileStruct, error) {

	// mengambil data user dari DB berdasarkan email
	query := "SELECT id, email, password, role FROM users WHERE email = $1"
	var result models.UsersStruct
	err := u.db.QueryRow(ctx, query, auth.Email).Scan(&result.Id, &result.Email, &result.Password, &result.Role)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.UsersStruct{}, models.ProfileStruct{}, err
	}

	queryGetProfile := "SELECT first_name, last_name, phone_number, photo_path, point FROM profiles WHERE user_id = $1"
	var profile models.ProfileStruct
	if err := u.db.QueryRow(ctx, queryGetProfile, result.Id).Scan(&profile.First_name, &profile.Last_name, &profile.Phone_number, &profile.PhotoPath, &profile.Point); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.UsersStruct{}, models.ProfileStruct{}, err
	}

	return result, profile, nil
}

// Repository get profile by id
func (u *UserRepository) GetProfileById(ctx context.Context, idInt int) ([]models.ProfileStruct, error) {

	// jika tidak terdapat data di redis maka jalankan query GET profile berikut ini
	query := "SELECT p.user_id, p.first_name, p.last_name, p.phone_number, p.photo_path, p.title, p.point, u.email FROM profiles p join users u on u.id = p.user_id WHERE user_id = $1"
	rows, err := u.db.Query(ctx, query, idInt)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	var result []models.ProfileStruct

	for rows.Next() {
		var profile models.ProfileStruct
		if err := rows.Scan(&profile.User_Id, &profile.First_name, &profile.Last_name, &profile.Phone_number, &profile.PhotoPath, &profile.Title, &profile.Point, &profile.Email); err != nil {
			return []models.ProfileStruct{}, err
		}
		result = append(result, profile)
	}
	
	return result, nil
}

// Repository update profile
func (u *UserRepository) UpdateProfile(ctx context.Context, idUser int, firstName, lastName, phoneNumber, title, password, email, passTrigger string) (pgconn.CommandTag, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	defer tx.Rollback(ctx)

	// update table profile berdasarkan user_id
	query := "UPDATE profiles SET first_name = $1, last_name = $2, phone_number = $3, title = $4, modified_at = $5 WHERE user_id = $6"
	values := []any{firstName, lastName, phoneNumber, title, time.Now(), idUser}
	cmd, err := tx.Exec(ctx, query, values...)
	if err != nil {
		log.Println("[ERROR : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	// update email berdasarkan user id
	queryNewEmail := "update users set email = $1 where id = $2"
	if _, err := tx.Exec(ctx, queryNewEmail, email, idUser); err != nil {
		log.Println("[ERROR : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	if passTrigger != "" {
		// melakukan update password bersadarkan user_id
		queryNewPassword := "update users set password = $1 where id = $2"
		if _, err := tx.Exec(ctx, queryNewPassword, password, idUser); err != nil {
			log.Println("[ERROR : ", err.Error())
			return pgconn.CommandTag{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR : ", err.Error())
		return pgconn.CommandTag{}, err
	}
	
	return cmd, nil
}

func (u *UserRepository) UpdatePhotoProfile(ctx context.Context, idUser int, filePath string) (pgconn.CommandTag, error) {
	newPath := strings.TrimPrefix(filePath, "public\\")
	// update table profile berdasarkan user_id
	query := "UPDATE profiles SET photo_path = $1, modified_at = $2 WHERE user_id = $3"
	values := []any{newPath, time.Now(), idUser}
	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		log.Println("[ERROR : ", err.Error())
		return pgconn.CommandTag{}, err
	}
	
	return cmd, nil
}