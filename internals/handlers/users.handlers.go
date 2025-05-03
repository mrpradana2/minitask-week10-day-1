package handlers

import (
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

// handler add user (fix)
func (u *UsersHandler) UserRegister(ctx *gin.Context) {
	// deklarasi body dari input user
	newDataUser := models.SignupPayload{}

	// binding data 
	// membaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
		log.Println("Binding error:", err)

		// error jika format email salah
		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: "failed",
				Msg: "incorrect email format",
			})
			return
		}

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: "failed",
				Msg: "password length must be at least 8 characters",
			})
			return
		}

		// error yang lainnya
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

	// error jika gagal mengkonversi password menjadi hash
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

	// lakukan error handling jika terjadi kesalahan dalam menjalankan query / user sudah terdaftar sebelumnya
	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusConflict, models.Message{
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

// handler user login (fix)
func (u *UsersHandler) UserLogin(ctx *gin.Context) {
	// mengambil body dari json / input user
	auth := models.UsersStruct{}

	// binding data 
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&auth); err != nil {
		log.Println(err.Error())

		// error jika format email salah
		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: "failed",
				Msg: "incorrect email format",
			})
			return
		}

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: "failed",
				Msg: "password length must be at least 8 characters",
			})
			return
		}

		// error yang lainnya
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "invalid data sent",
		})
		return
	}

	// melakukan eksekusi fungsi repository user login
	result, err := u.usersRepo.UserLogin(ctx.Request.Context(), auth)

	// error jika terjadi kesalahan dalam mengakses server
	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	// mengecek apakan input password dari user sama dengan hasil pencarian user di database
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

	// error jika password tidak sesuai dengan password yang di hash
	if !valid {
		ctx.JSON(http.StatusUnauthorized, models.Message{
			Status: "failed",
			Msg: "incorrect email or password",
		})
		return
	}

	// jika berhasil login, maka berikan identitas (jwt)
	claims := pkg.NewClaims(result.Id, result.Role)
	token, err := claims.GenerateToken()

	// error jika gagal menghasilkan token
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	} 

	// pesan jika user berhasil login
	ctx.JSON(http.StatusOK, models.Message{
		Status: "success",
		Msg: "login success",
		Token: token,
	})
}

// handler get profile (fix)
func (u *UsersHandler) GetProfileById(ctx *gin.Context) {
	// ambil data user yang login
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// eksekusi fungsi repository get profile
	result, err := u.usersRepo.GetProfileById(ctx.Request.Context(), userClaims.Id)

	// error jika profile user tidak ditemukan
	if err != nil {
		log.Println("[ERROR] :", err)
		ctx.JSON(http.StatusNotFound, models.Message{
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

// handler update profile (fix)
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
	// sediakan variabel untuk menapung input dari form
	var formBody models.ProfileStruct

	// error jika data yang diinput tidak sesuai
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Println(err.Error())

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'NewPassword'") {
			ctx.JSON(http.StatusBadRequest, models.Message{
				Status: "failed",
				Msg: "password length must be at least 8 characters",
			})
			return
		}

		// jika terjadi error lain 
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan server",
		})
		return
	}

	// ambil nilai dari forn yang dikirim user
	file := formBody.Photo_path
	firstName := formBody.First_name
	lastName := formBody.Last_name
	phoneNumber := formBody.Phone_number
	title := formBody.Title
	newPassword := formBody.NewPassword
	confirmPassword := formBody.ComfirmPassword

	// error jika newpassword dan confirm password tidak sama
	if newPassword != confirmPassword {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
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
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "hash failed",
		})
		return
	}

	// error = jika file tidak diupload user
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
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
	filepath := fp.Join("public", "img", filename)

	// error jika gagal melakukan upload
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan upload",
		})
		return
	}

	// jalankan fungsi repository untuk update data
	cmd, err := u.usersRepo.UpdateProfile(ctx.Request.Context(), idInt, firstName, lastName, phoneNumber, filepath, title, hashedPass)

	// error jika gagal mengakses server
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