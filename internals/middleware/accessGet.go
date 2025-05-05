package middleware

import (
	"log"
	"net/http"
	"slices"
	"tikcitz-app/internals/models"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AcceessGate(allowedRole ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// 1. ambil payload/claims dari context gin
		claims, exist := ctx.Get("Payload")
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
				Status: "failed",
				Msg: "silahkan login terlebih dahulu",
			})
			return
		} 
		
		// type assertion claims menjadi pkg.claims
		userClaims, ok := claims.(*pkg.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
				Status: "failed",
				Msg: "identitas login anda rusak, silahkan login kembali",
			})
			return
		}
		log.Println("ALLOWED ROLE ", allowedRole)
		log.Println("USER CLAIMS ", userClaims.Role)
		// cek role yang ada di claims
		if !slices.Contains(allowedRole, userClaims.Role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.Message{
				Status: "failed",
				Msg: "Anda tidak dapat mengakses sumber ini",
			})
			return
		}
		ctx.Next()
	}
}

func (m *Middleware) AcceessGateAdmin(ctx *gin.Context) {
	// 1. ambil payload/claims dari context gin
	claims, exist := ctx.Get("Payload")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
			Status: "failed",
			Msg: "silahkan login terlebih dahulu",
		})
		return
	} 
	
	// type assertion claims menjadi pkg.claims
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Message{
			Status: "failed",
			Msg: "identitas login anda rusak, silahkan login kembali",
		})
		return
	}

	// cek role yang ada di claims
	if userClaims.Role != "admin" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, models.Message{
			Status: "failed",
			Msg: "Anda tidak dapat mengakses sumber ini",
		})
		return
	}
	ctx.Next()
}