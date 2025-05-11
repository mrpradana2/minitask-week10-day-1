package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/middleware"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterSchedule(router *gin.Engine, scheduleRepo *repositories.ScheduleRepository, middleware *middleware.Middleware) {
	routerSchedule := router.Group("/schedule")

	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo)
	
	// router get schedule
	routerSchedule.GET("/:movieId", middleware.VerifyToken, middleware.AcceessGate("admin", "user"), scheduleHandler.GetScheduleMovie)
}