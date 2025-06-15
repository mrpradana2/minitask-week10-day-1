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

type SeatsHandler struct{
	seatsRepo *repositories.SeatsRepository
}
func NewSeatshandler(seatsRepo *repositories.SeatsRepository) *SeatsHandler {
	return &SeatsHandler{seatsRepo: seatsRepo}
} 

// Get Available Seat
// @summary					Get available seat
// @Description 			get available seat by schedule id
// @Tags        			Seats
// @router					/seats/:scheduleId [get]
// @Param        			scheduleId query string true "schedule id to get available seat"
// @Accept					json
// @produce					json
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					404 {object} models.MessageNotFound
// @success					200 {object} models.MessageResult
func (s *SeatsHandler) GetSeatsAvailable(ctx *gin.Context) {

	// mendapatkan id dari params
	idStr, ok := ctx.Params.Get("scheduleId")

	// handling error jika param tidak ada
	if !ok {
		log.Println("[ERROR] : ", errors.New("params not found"))
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "Param id is needed",
		})
		return
	}

	// melakukan konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	// melakukan error handling jika gagal mengkonversi id
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// menjalankan fungsi repository getseatAvailable untuk mengakses data dari server
	seats, err := s.seatsRepo.GetSeatsAvailable(ctx.Request.Context(), models.SeatsStruct{}, idInt)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	// tampilkan hasilnya jika berhasil mendapatkan request dari server
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "success",
		Result: seats,
	})
}