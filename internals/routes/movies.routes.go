package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/middleware"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterMovies(router *gin.Engine, moviesRepo *repositories.MoviesRepository, middleware *middleware.Middleware) {
	routerMovies := router.Group("/movies")
	moviesHandler := handlers.NewMovieshandler(moviesRepo)

	// Router add movie (fix)
	routerMovies.POST("", middleware.VerifyToken, middleware.AcceessGate("admin", "user"), moviesHandler.AddMovie)

	// Router update movie
	routerMovies.PUT("/:id", middleware.VerifyToken, middleware.AcceessGate("admin", "user"), moviesHandler.UpdateMovie)

	// Router delete movie
	routerMovies.DELETE("/:id",  middleware.VerifyToken, middleware.AcceessGate("admin", "user"), moviesHandler.DeleteMovie)

	// router get movie upcoming
	routerMovies.GET("/moviesupcoming", moviesHandler.GetMoviesUpcoming)

	// router get movie popular
	routerMovies.GET("/moviespopular", moviesHandler.GetMoviesPopular)

	// router get detail movie
	routerMovies.GET("/:id", moviesHandler.GetDetailMovie)

	// router get movie with pagination
	routerMovies.GET("", moviesHandler.GetMoviesWithPagination)
}