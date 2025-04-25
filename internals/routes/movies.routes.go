package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterMovies(router *gin.Engine, moviesRepo *repositories.MoviesRepository) {
	// DATA MOVIES
	routerMovies := router.Group("/movies")
	moviesHandler := handlers.NewMovieshandler(moviesRepo)
	
	// get all movie
	routerMovies.GET("", moviesHandler.GetMovies)

	// add movie
	routerMovies.POST("/add", moviesHandler.AddMovie)
}