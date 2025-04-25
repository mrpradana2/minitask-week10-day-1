package handlers

import (
	"log"
	"net/http"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

type Movieshandler struct{}

func NewMovieshandler() *Movieshandler {
	return &Movieshandler{}
}

// Get all movies
func (m *Movieshandler) GetMovies(ctx *gin.Context) {
	result, err := repositories.MovieRepo.GetMovies(ctx)

	if err != nil {
		log.Println("[ERROR]", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "terjadi kesalahan sistem",
		})
		return
	}

	if len(result) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": result,
	})
}
