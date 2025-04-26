package handlers

import (
	"log"
	"net/http"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleRepo *repositories.ScheduleRepository
}

func NewScheduleHandler(scheduleRepo *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{scheduleRepo: scheduleRepo}
}

// handler get schedule
func (s *ScheduleHandler) GetScheduleMovie(ctx *gin.Context) {
	result, err := s.scheduleRepo.GetScheduleMovie(ctx, &models.ScheduleStruct{})
	if err != nil {
		log.Println("[ERROR]: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "schedule not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
	})
}