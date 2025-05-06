package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterSchedule(router *gin.Engine, scheduleRepo *repositories.ScheduleRepository) {
	routerSchedule := router.Group("/schedule")

	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)
	
	// router get schedule
	routerSchedule.GET("/:movieId", scheduleHandler.GetScheduleMovie)
}