package routes

import (
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

    moviesRepo := repositories.NewMoviesRepository(db)
    usersRepo := repositories.NewUserRepository(db)

	InitRouterUsers(router, usersRepo)
	InitRouterMovies(router, moviesRepo)

	return router

    // router.GET("movies/movie-upcoming", func(ctx *gin.Context) {
    //     result := []moviesStruct{}

    //     if len(movies) == 0 {
    //         ctx.JSON(http.StatusInternalServerError, gin.H{
    //             "msg": "an error occurred on the server",
    //         })
    //     }

    //     for _, movie := range movies {
    //         if movie.Status_movie == "movie upcoming" {
    //             result = append(result, movie)
    //         }
    //     }

    //     if len(result) == 0 {
    //         ctx.JSON(http.StatusNotFound, gin.H{
    //             "msg": "movie upcoming not found",
    //         })
    //     }

    //     ctx.JSON(http.StatusOK, gin.H{
    //         "msg": "success",
    //         "data": result,
    //     })
    // })
    
    // router.GET("movies/movie-popular", func(ctx *gin.Context) {
    //     result := []moviesStruct{}

    //     if len(movies) == 0 {
    //         ctx.JSON(http.StatusInternalServerError, gin.H{
    //             "msg": "an error occurred on the server",
    //         })
    //     }

    //     for _, movie := range movies {
    //         if movie.Status_movie == "movie popular" {
    //             result = append(result, movie)
    //         }
    //     }

    //     if len(result) == 0 {
    //         ctx.JSON(http.StatusNotFound, gin.H{
    //             "msg": "movie popular not found",
    //         })
    //     }

    //     ctx.JSON(http.StatusOK, gin.H{
    //         "msg": "success",
    //         "data": result,
    //     })
    // })
}