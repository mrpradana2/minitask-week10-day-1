package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	fp "path/filepath"
	"strings"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"
	"tikcitz-app/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct{
	usersRepo *repositories.UserRepository
}

func NewUsersHandlers(usersRepo *repositories.UserRepository) *UsersHandler {
	return &UsersHandler{usersRepo: usersRepo}
}


// Register
// @summary					Register user
// @router					/users/signup [post]
// @Description 			Register with email and password for access application 
// @Tags        			Users
// @Param        			login body models.UserLogin true "Input email and password"
// @accept					json
// @produce					json
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					409 {object} models.MessageConflict
// @success					201 {object} models.MessageCreated
func (u *UsersHandler) UserRegister(ctx *gin.Context) {
	// deklarasi body dari input user
	newDataUser := models.UsersStruct{}

	// binding data 
	// membaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
		log.Println("[ERROR] :", err.Error())

		// error jika format email salah
		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: http.StatusBadRequest,
				Msg: "incorrect email format",
			})
			return
		}

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: http.StatusBadRequest,
				Msg: "password length must be at least 8 characters",
			})
			return
		}

		// error yang lainnya
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "invalid data sent",
		})
		return
	}

	// hash password, mengkonversi password normal menjadi bentuk lain yang sulit dibaca
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenHashedPassword(newDataUser.Password)

	// error jika gagal mengkonversi password menjadi hash
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "hash failed",
		})
		return
	}

	// default value untuk role
	role := "user"

	// eksekusi fungsi repository register user
	cmd, err := u.usersRepo.UserRegister(ctx.Request.Context(), newDataUser.Email, hashedPass, role)

	// lakukan error handling jika terjadi kesalahan dalam menjalankan query / user sudah terdaftar sebelumnya
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusConflict, models.Message{
			Status: http.StatusConflict,
			Msg: "user already registered",
		})
		return
	}

	// cek apakah perintah berhasil manambahkan data di database 
	if cmd.RowsAffected() == 0 {
		log.Println("[ERROR] : ", errors.New("query failed, did not change the data in the database"))
		log.Println("Query failed, did not change the data in the database")
		return
	}

	// return jika server berhasil memberikan response
	ctx.JSON(http.StatusCreated, models.Message{
		Status: http.StatusCreated,
		Msg: "successfully create an account",
	})
}


// Login
// @summary					Login user
// @router					/users/login [post]
// @Description 			Request data login email and password to authentication login
// @Tags        			Users
// @Param        			login body models.UserLogin true "Input email and password"
// @accept					json
// @produce					json
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					401 {object} models.MessageUnauthorized
// @success					200 {object} models.MessageLogin
func (u *UsersHandler) UserLogin(ctx *gin.Context) {
	// mengambil body dari json / input user
	auth := models.UsersStruct{}

	// binding data 
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&auth); err != nil {
		log.Println("[ERROR] : ", err.Error())

		// error jika format email salah
		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: http.StatusBadRequest,
				Msg: "incorrect email format",
			})
			return
		}

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: http.StatusBadRequest,
				Msg: "password length must be at least 8 characters",
			})
			return
		}

		// error yang lainnya
		log.Println("[ERROR] : ", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "invalid data sent",
		})
		return
	}

	// melakukan eksekusi fungsi repository user login
	result, profile, err := u.usersRepo.UserLogin(ctx.Request.Context(), auth)
	// error jika terjadi kesalahan dalam mengakses server
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusUnauthorized, models.Message{
			Status: http.StatusUnauthorized,
			Msg: "incorrect email or password",
		})
		return
	}

	// mengecek apakan input password dari user sama dengan hasil pencarian user di database
	hash := pkg.InitHashConfig()
	valid, err := hash.CompareHashAndPassword(result.Password, auth.Password)

	if err != nil {
		log.Println("[ERROR] : ",err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// error jika password tidak sesuai dengan password yang di hash
	if !valid {
		log.Println("[ERROR] : ", errors.New("incorrect email or password"))
		ctx.JSON(http.StatusUnauthorized, models.Message{
			Status: http.StatusUnauthorized,
			Msg: "incorrect email or password",
		})
		return
	}

	// jika berhasil login, maka berikan identitas (jwt)
	claims := pkg.NewClaims(result.Id, result.Role)

	token, err := claims.GenerateToken()

	// error jika gagal menghasilkan token
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	} 

	// pesan jika user berhasil login
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "login success",
		Token: token,
		Result: profile,
	})
}


// GetProfile handles the request to fetch user profile
// @Summary     			Get profile
// @Description 			Get data profile user
// @Router      			/users [get]
// @Tags        			Users
// @Param 					Authorization header string true "Bearer Token"
// @Accept      			json
// @Produce     			json
// @Security    			BearerAuth
// @Failure     			500 {object} models.MessageInternalServerError
// @Failure     			404 {object} models.MessageNotFound
// @Failure     			401 {object} models.MessageUnauthorized
// @Success     			200 {object} models.MessageResult
func (u *UsersHandler) GetProfileById(ctx *gin.Context) {
	// ambil data user yang login
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// eksekusi fungsi repository get profile
	result, err := u.usersRepo.GetProfileById(ctx.Request.Context(), userClaims.Id)

	// error jika profile user tidak ditemukan
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "user not found",
		})
		return
	}

	// pesan json jika berhasil mendapatkan profile user
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "success get profile",
		Result: result,
	})
}

// Update Profile
// @summary					Update profile
// @router					/users [patch]
// @Description 			Upload a new data profile for the user
// @Tags        			Users
// @Param        			requestBody body models.RequestProfileStruct true "Update profile request"
// @Param 					Authorization header string true "Bearer Token"
// @accept					json
// @produce					json
// @Security    			BearerAuth
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					401 {object} models.MessageUnauthorized
// @failure					404 {object} models.MessageNotFound
// @success					200 {object} models.MessageOK
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {

	updateProfile := models.ProfileStruct{}

	if err := ctx.ShouldBindJSON(&updateProfile); err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "invalid data sent",
		})
		return
	}

	firstName := updateProfile.First_name
	lastName := updateProfile.Last_name
	phoneNumber := updateProfile.Phone_number
	email := updateProfile.Email
	title := updateProfile.Title
	newPassword := updateProfile.NewPassword
	confirmPassword := updateProfile.ConfirmPassword

	// error jika newpassword dan confirm password tidak sama
	if newPassword != confirmPassword {
		log.Println("[ERROR] : ", errors.New("password are not the same"))
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "password are not the same",
		})
		return
	}

	// hash password, mengkonversi password normal menjadi bentuk lain yang sulit dibaca
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenHashedPassword(confirmPassword)

	// error jika gagal mengkonversi password menjadi hash
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "hash failed",
		})
		return
	}

	// ambil data user yang login
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// ambil id user yang login
	idInt := userClaims.Id

	// jalankan fungsi repository untuk update data
	cmd, err := u.usersRepo.UpdateProfile(ctx.Request.Context(), idInt, firstName, lastName, phoneNumber, title, hashedPass, email, confirmPassword)

	// error jika gagal mengakses server
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// pesan error jika tidak terjadi perubahan pada database atau jika user id tidak ditemukan
	if cmd.RowsAffected() == 0 {
		log.Println("[ERROR] : ", errors.New("query failed, did not change the data in the database"))
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "no updated data",
		})
		return
	}

	// menampilkan hasil jika berhasil mengupdate profile
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "update success",
	})
}

// Update Photo Profile
// @summary					Update photo profile
// @router					/users/photoProfile [patch]
// @Description 			Upload a new profile photo for the user
// @Tags        			Users
// @Param        			requestBody body models.RequestPhotoProfileStruct true "Update photo profile request"
// @Param 					Authorization header string true "Bearer Token"
// @accept					multipart/form-data
// @produce					json
// @Security    			BearerAuth
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					401 {object} models.MessageUnauthorized
// @failure					404 {object} models.MessageNotFound
// @success					200 {object} models.MessageOK
func (u *UsersHandler) UpdatePhotoProfile(ctx *gin.Context) {
	// sediakan variabel untuk menapung input dari form
	var formBody models.PhotoProfileStruct

	// error jika data yang diinput tidak sesuai
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Println("[ERROR] : ", err.Error()) 
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "terjadi kesalahan server",
		})
		return
	}

	// ambil nilai dari forn yang dikirim user
	file := formBody.Photo_path

	// error = jika file tidak diupload user
	if file == nil {
		log.Println("[ERROR] : ", errors.New("file not found"))
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "your file is empty",
		})
		return
	}

	// ambil data user yang login
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// ambil id user yang login
	idInt := userClaims.Id

	// ambil extensi file
	ext := fp.Ext(file.Filename)

	// rakit file name sebagai nama dari file yang akan disimpan
	filename := fmt.Sprintf("%d_%d_user_image%s", time.Now().UnixNano(), userClaims.Id, ext)

	// gabungkan alamat folder dengan nama file yang sidah dibuat sehingga menjadi filepath
	filepath := fp.Join("public", "img", "photoProfiles", filename)

	// error jika gagal melakukan upload
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "terjadi kesalahan upload",
		})
		return
	}

	// jalankan fungsi repository untuk update data
	cmd, err := u.usersRepo.UpdatePhotoProfile(ctx.Request.Context(), idInt, filepath)

	// error jika gagal mengakses server
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// pesan error jika tidak terjadi perubahan pada database atau jika user id tidak ditemukan
	if cmd.RowsAffected() == 0 {
		log.Println("[ERROR] : ", errors.New("query failed, did not change the data in the database"))
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "no updated data",
		})
		return
	}

	// menampilkan hasil jika berhasil mengupdate profile
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "update success",
	})
}

// // // handler user auth
// func (u *UsersHandler) VerifyToken(ctx *gin.Context) {
// 	// 1. ambil token dari header
// 	bearerToken := ctx.GetHeader("Authorization")

// 	// cek jika bearerToken kosong
// 	if bearerToken == "" {
// 		ctx.JSON(http.StatusUnauthorized, models.Message{
// 			Status: "failed",
// 			Msg: "silahkan login terlebih dahulu",
// 		})
// 		return
// 	}

// 	// 2. pisahkan token dari bearer
// 	token := strings.Split(bearerToken, " ")[1]

// 	if token == "" {
// 		ctx.JSON(http.StatusUnauthorized, models.Message{
// 			Status: "failed",
// 			Msg: "silahkan login terlebih dahulu",
// 		})
// 		return
// 	}

// 	// verifikasi token
// 	claims := &pkg.Claims{}
// 	if err := claims.VerifyToken(token); err != nil {
// 		log.Println(err.Error())
// 		if err.Error() == "expired token" || err.Error() == "token has invalid claims: token is expired" {
// 			ctx.JSON(http.StatusUnauthorized, models.Message{
// 				Status: "failed",
// 				Msg: "silahkan login kembali",
// 			})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, models.Message{
// 			Status: "failed",
// 			Msg: "terjadi kesalahan server",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.Message{
// 		Status: "ok",
// 		Msg: "success",
// 		Result: claims,
// 	})
// }