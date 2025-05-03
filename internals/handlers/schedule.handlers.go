package handlers

import (
	"log"
	"net/http"
	"strconv"
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
	// mengambil user id di params
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "failed",
			Msg: "params id is needed",
		})
		return
	}

	// konversi id string menjadi integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[ERROR] : ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "an error occurred on the server",
		})
		return
	}

	result, err := s.scheduleRepo.GetScheduleMovie(ctx, &models.ScheduleStruct{}, idInt)
	if err != nil {
		log.Println("[ERROR]: ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "ok",
			Msg: "an error occurred on the server",
		})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "failed",
			Msg: "schedule not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
}