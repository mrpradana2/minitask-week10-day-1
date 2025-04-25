package handlers

import (
	"log"
	"net/http"
	"strconv"
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

// Handler get all movies
func (m *Movieshandler) GetMovies(ctx *gin.Context) {
	result, err := m.moviesRepo.GetMovies(ctx.Request.Context())

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

// Handler add movie
func (m *Movieshandler) AddMovie(ctx *gin.Context)  {
	newDataMovie := models.MoviesStruct{}

	if err := ctx.ShouldBindJSON(&newDataMovie); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	cmd, err := m.moviesRepo.AddMovie(ctx.Request.Context(), newDataMovie)

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

// Handler update movie
func (m *Movieshandler) UpdateMovie(ctx *gin.Context) {
	var updateMovie models.MoviesStruct

	// binding data
	if err := ctx.ShouldBindJSON(&updateMovie); err != nil {
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

	cmd, err := m.moviesRepo.UpdateMovie(ctx.Request.Context(), &updateMovie, idInt)
	if err != nil {
		log.Println("Insert profile error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, could not change the data in the database")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "data successfully changed",
	})
}

// handler delete movie
func (m *Movieshandler) DeleteMovie(ctx *gin.Context) {
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

	cmd, err := m.moviesRepo.DeleteMovie(ctx.Request.Context(), idInt)

	if err != nil {
		log.Println("Insert profile error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if cmd.RowsAffected() == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "data successfully deleted",
	})
}

// handler get movie upcoming
func (m *Movieshandler) GetMoviesUpcoming(ctx *gin.Context) {
	result, err := m.moviesRepo.GetMovieUpcoming(ctx.Request.Context())
	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "movie not found",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
	})
}
