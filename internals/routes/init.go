package routes

import (
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

    repositories.NewMoviesRepository()
    repositories.NewUserRepository()

	InitRouterUsers(router)
	InitRouterMovies(router)

	return router

    // router.POST("movies", func(ctx *gin.Context) {
    //     newMovie := &moviesStruct{}
    //     if err := ctx.ShouldBind(newMovie); err != nil {
    //         ctx.JSON(http.StatusInternalServerError, gin.H{
    //             "msg": "terjadi kesalahan sisten",
    //         })
    //         return
    //     }

    //     newMovies := append(movies, *newMovie)
    //     ctx.JSON(http.StatusOK, gin.H{
    //         "msg": "success",
    //         "data": newMovies,
    //     })
    // })

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