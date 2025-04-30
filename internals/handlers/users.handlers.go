package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct{
	usersRepo *repositories.UserRepository
}

func NewUsersHandlers(usersRepo *repositories.UserRepository) *UsersHandler {
	return &UsersHandler{usersRepo: usersRepo}
}

// handler add user
func (u *UsersHandler) UserRegister(ctx *gin.Context) {
	// delkarasi body dari input user
	newDataUser := models.SignupPayload{}

	// binding data 
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "invalid data sent",
		})
		return
	}

	// hash password, mengkonversi password normal menjadi bentuk lain yang sulit dibaca
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenHashedPassword(newDataUser.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "hash failed",
		})
		return
	}

	// default value untuk role
	role := "user"

	// eksekusi fungsi repository register user
	cmd, err := u.usersRepo.UserRegister(ctx.Request.Context(), newDataUser.Email, hashedPass, role)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "user already registered",
		})
		return
	}

	// cek apakah perintah berhasil manambahkan data di database 
	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, did not change the data in the database")
		return
	}

	// return jika server berhasil memberikan response
	ctx.JSON(http.StatusCreated, models.Message{
		Status: "success",
		Msg: "successfully create an account",
	})
}

// handler user login
func (u *UsersHandler) UserLogin(ctx *gin.Context) {
	// mengambil body dari json / input user
	auth := models.UsersStruct{}

	// binding data 
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&auth); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "invalid data sent",
		})
		return
	}

	// melakukan eksekusi fungsi repository user login
	result, err := u.usersRepo.UserLogin(ctx.Request.Context(), auth)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// mengecek apakan input dari user sama dengan hasil pencarian user di database
	hash := pkg.InitHashConfig()
	valid, err := hash.CompareHashAndPassword(result.Password, auth.Password)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// jika pengecekan password tidak sesuai
	if !valid {
		ctx.JSON(http.StatusUnauthorized, models.Message{
			Status: "failed",
			Msg: "incorrect username or password",
		})
		return
	}

	// jika user berhasil login
	ctx.JSON(http.StatusOK, models.Message{
		Status: "success",
		Msg: "login success",
	})
}

// handler get profile
func (u *UsersHandler) GetProfileById(ctx *gin.Context) {
	// mengambil user id di params
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "params id is needed",
		})
		return
	}

	// konversi id string menjadi integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[ERROR] : ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// eksekusi fungsi repository get profile
	result, err := u.usersRepo.GetProfileById(ctx.Request.Context(), idInt)

	// error jika profile user tidak ditemukan
	if err != nil {
		log.Println("[ERROR] :", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "user not found",
		})
		return
	}

	// pesan json jika berhasil mendapatkan profile user
	ctx.JSON(http.StatusOK, models.Message{
		Status: "success",
		Msg: "success get profile",
		Result: result,
	})
}

// handler update profile
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
	// mengambil data dari body json / input user 
	var updateProfile models.ProfileStruct

	// binding data
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error
	if err := ctx.ShouldBindJSON(&updateProfile); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "invalid data sent",
		})
		return
	}

	// ambil parameter berdasarkan user id
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "param id is needed",
		})
		return
	}

	// konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// jalankan fungsi repository untuk update data
	cmd, err := u.usersRepo.UpdateProfile(ctx.Request.Context(), updateProfile, idInt)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// pesan error jika tidak terjadi perubahan pada database atau jika user id tidak ditemukan
	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, did not change the data in the database")
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "failed",
			Msg: "no updated data",
		})
		return
	}

	// menampilkan hasil jika berhasil mengupdate profile
	ctx.JSON(http.StatusOK, models.Message{
		Status: "success",
		Msg: "update success",
	})
}
