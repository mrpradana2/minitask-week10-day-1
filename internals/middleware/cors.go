package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) CORSMiddleware(ctx *gin.Context) {
	// setup whitlist origin
	whiteListOrigin := []string{"http://localhost:5173", "http://localhost:3000", "http://localhost"}
	origin := ctx.GetHeader("Origin")

	// allowedOrigin := ""

	log.Println("[DEBUG] Origin: ", origin)
	if slices.Contains(whiteListOrigin, origin) {
		log.Println("[DEBUG] whitlisted")
		ctx.Header("Access-Control-Allow-Origin", origin)
		// allowedOrigin = origin
	}

	// ctx.Header("Access-Control-Allow-Origin", allowedOrigin)
	ctx.Header("Access-Control-Allow-MethodS", "GET, POST, HEAD, PATCH, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, my-headers")

	// handle preflight // or "OPTIONS"
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}