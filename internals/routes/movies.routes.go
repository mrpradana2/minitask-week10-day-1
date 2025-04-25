package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterMovies(router *gin.Engine, moviesRepo *repositories.MoviesRepository) {
	routerMovies := router.Group("/movies")
	moviesHandler := handlers.NewMovieshandler(moviesRepo)
	
	// Router get all movie
	routerMovies.GET("", moviesHandler.GetMovies)

	// Router add movie
	routerMovies.POST("/add", moviesHandler.AddMovie)

	// Router update movie
	routerMovies.PUT("/edit/:id", moviesHandler.UpdateMovie)
}