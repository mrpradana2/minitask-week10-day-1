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
			Status: "failes",
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
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Param id is needed",
		})
		return
	}

	// konversi id string menjadi integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan server",
		})
		return
	}

	// eksekusi fungsi repository 
	result, err := u.usersRepo.GetProfileById(ctx.Request.Context(), idInt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "user profile not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": result,
	})
}

// handler update profile
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
	var updateProfile models.ProfileStruct

	// binding data
	if err := ctx.ShouldBindJSON(&updateProfile); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Param id is needed",
		})
		return
	}

	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
		return
	}

	cmd, err := u.usersRepo.UpdateProfile(ctx.Request.Context(), updateProfile, idInt)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, did not change the data in the database")
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "update success",
	})
}
