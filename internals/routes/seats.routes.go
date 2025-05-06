package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterSeats(router *gin.Engine, seatsRepo *repositories.SeatsRepository) {
	routerSeats := router.Group("/seats")
	seatsHandler := handlers.NewSeatshandler(seatsRepo)

	// router get seats avaliable
	routerSeats.GET("/:scheduleId", seatsHandler.GetSeatsAvailable)
}