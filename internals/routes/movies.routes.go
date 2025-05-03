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

	// middleware := middleware.InitMiddleware()
	
	// Router get all movie (fix)
	routerMovies.GET("", moviesHandler.GetMovies)

	// Router add movie (fix)
	routerMovies.POST("/add", middleware.VerifyToken, middleware.AcceessGate("admin"), moviesHandler.AddMovie)

	// Router update movie
	routerMovies.PATCH("/edit/:id", middleware.VerifyToken, middleware.AcceessGate("admin"), moviesHandler.UpdateMovie)

	// Router delete movie
	routerMovies.DELETE("/delete/:id", moviesHandler.DeleteMovie)

	// router get movie upcoming
	routerMovies.GET("/moviesupcoming", moviesHandler.GetMoviesUpcoming)

	// router get movie popular
	routerMovies.GET("/moviespopular", moviesHandler.GetMoviesPopular)

	// router get detail movie
	routerMovies.GET("/detail/:id", moviesHandler.GetDetailMovie)

	// router get movie with pagination
	routerMovies.GET("/page", moviesHandler.GetMoviesWithPagination)
}