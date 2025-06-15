package middleware

import (
	"log"
	"net/http"
	"strings"
	"tikcitz-app/internals/models"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func InitMiddleware() *Middleware {
	return &Middleware{}
}

func (u *Middleware) VerifyToken(ctx *gin.Context) {
	// 1. ambil token dari header
	bearerToken := ctx.GetHeader("Authorization")
	// cek jika bearerToken kosong
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
			Status: http.StatusUnauthorized,
			Msg:    "silahkan login terlebih dahuluu",
		})
		return
	}

	// verifikasi bearer token
	if !strings.Contains(bearerToken, "Bearer") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
			Status: http.StatusUnauthorized,
			Msg: "silahkan login terlebih dahulu",
		})
		return
	}

	// 2. pisahkan token dari bearer
	token := strings.Split(bearerToken, " ")[1]

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
			Status: http.StatusUnauthorized,
			Msg:    "silahkan login terlebih dahulu",
		})
		return
	}

	// verifikasi token
	claims := &pkg.Claims{}
	if err := claims.VerifyToken(token); err != nil {
		log.Println(err.Error())
		// if err.Error() == "expired token" || err.Error() == "token has invalid claims: token is expired" {
		if strings.Contains(err.Error(), "expired") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
				Status: http.StatusUnauthorized,
				Msg:    "sesi anda berakhir, silahkan login kembali",
			})
			return
		}
		// format salah
		if strings.Contains(err.Error(), "malformed") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
				Status: http.StatusUnauthorized,
				Msg:    "identitas login anda rusak, silahkan login kembali",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusUnauthorized,
			Msg:    "terjadi kesalahan server",
		})
		return
	}

	// ctx.JSON(http.StatusOK, models.Message{
	// 	Status: "ok",
	// 	Msg:    "success",
	// 	Result: claims,
	// })

	ctx.Set("Payload", claims)
	ctx.Next()
}