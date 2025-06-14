package routes

import (
	"tikcitz-app/internals/middleware"
	"tikcitz-app/internals/repositories"

	_ "tikcitz-app/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool,rdb *redis.Client) *gin.Engine {
	router := gin.Default()
    middleware := middleware.InitMiddleware()

    moviesRepo := repositories.NewMoviesRepository(db, rdb)
    usersRepo := repositories.NewUserRepository(db, rdb)
    scheduleRepo := repositories.NewScheduleRepository(db)
    ordersRepo := repositories.NewOrdersRepository(db)
    seatsRepo := repositories.NewSeatsRepository(db)
    router.Use(middleware.CORSMiddleware)

    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))

	InitRouterUsers(router, usersRepo, middleware)
	InitRouterMovies(router, moviesRepo, middleware)
    InitRouterSchedule(router, scheduleRepo, middleware)
    InitRouterOrders(router, ordersRepo, middleware)
    InitRouterSeats(router, seatsRepo)

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