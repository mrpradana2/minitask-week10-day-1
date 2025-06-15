package handlers

import (
	"errors"
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

// Get Schedule Movie
// @summary					Get Schedule movie
// @router					/schedule/:movieId [get]
// @Description 			Get schedule movie by movie id
// @Tags        			Schedule
// @Param        			movieId query string true "get schedule by movie id"
// @Accept					json
// @produce					json
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					404 {object} models.MessageNotFound
// @success					200 {object} models.MessageResult
func (s *ScheduleHandler) GetScheduleMovie(ctx *gin.Context) {
	// mengambil user id di params
	idStr, ok := ctx.Params.Get("movieId")

	// handling error jika param tidak ada
	if !ok {
		log.Println("[ERROR] : ", errors.New("params not found"))
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "params id is needed",
		})
		return
	}

	// konversi id string menjadi integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// menjalankan fungsi repository getschedulemovie
	result, err := s.scheduleRepo.GetScheduleMovie(ctx, &models.ScheduleStruct{}, idInt)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// error handling jika data yang diambil dari server kosong
	if len(result) < 1 {
		log.Println("[ERROR] : ", errors.New("schedule not found"))
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "schedule not found",
		})
		return
	}

	// menambahkan id di result

	// tampilkan hasil response dari server
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "success",
		Result: result,
	})
}