package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"

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
	newDataUser := models.SignupPayload{}

	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	cmd, err := u.usersRepo.UserRegister(ctx.Request.Context(), newDataUser)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occured on the server"})
	}

	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, did not change the data in the database")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

// handler user login
func (u *UsersHandler) UserLogin(ctx *gin.Context) {
	auth := models.UsersStruct{}

	if err := ctx.ShouldBindJSON(&auth); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	value := []any{auth.Email, auth.Password}

	result, err := u.usersRepo.UserLogin(ctx.Request.Context(), auth)

	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	// mengecek apakan input dari user sama dengan hasil pencarian user di database
	var userLogin []models.UsersStruct
	if value[0] == result.Email && value[1] == result.Password {
		userLogin = append(userLogin, result)
	}

	if len(userLogin) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "incorrect email or password",
		})
		return
	}

	// jika user berhasil login
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "login user success",
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

	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan server",
		})
		return
	}

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
