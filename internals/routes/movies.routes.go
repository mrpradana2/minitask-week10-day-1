package routes

import (
	"tikcitz-app/internals/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouterMovies(router *gin.Engine) {
	// DATA MOVIES
	routerMovies := router.Group("/movies")
	moviesHandler := handlers.NewMovieshandler()
	
	// get all movie
	routerMovies.GET("", moviesHandler.GetMovies)

	// add movie
	routerMovies.POST("/add", moviesHandler.AddMovie)
}