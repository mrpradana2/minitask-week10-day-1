package handlers

import (
	"log"
	"net/http"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

type Movieshandler struct{
	moviesRepo *repositories.MoviesRepository
}

func NewMovieshandler(moviesRepo *repositories.MoviesRepository) *Movieshandler {
	return &Movieshandler{moviesRepo: moviesRepo}
}

// Get all movies
func (m *Movieshandler) GetMovies(ctx *gin.Context) {
	result, err := m.moviesRepo.GetMovies(ctx)

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

// Add movie
func (m *Movieshandler) AddMovie(ctx *gin.Context)  {
	newDataMovie := models.MoviesStruct{}

	if err := ctx.ShouldBindJSON(&newDataMovie); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	cmd, err := m.moviesRepo.AddMovie(ctx, newDataMovie)

	if err != nil {
		log.Println("Insert profile error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Terjadi kesalahan server"})
	}

	if cmd.RowsAffected() == 0 {
		log.Println("Query Gagal, Tidak merubah data di DB")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
