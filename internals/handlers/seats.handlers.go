package handlers

import (
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

func (s *SeatsHandler) GetSeatsAvailable(ctx *gin.Context) {
	// mendapatkan id dari params
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "Param id is needed",
		})
		return
	}

	// melakukan konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	// melakukan error handling jika gagal mengkonversi id
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	seats, err := s.seatsRepo.GetSeatsAvailable(ctx.Request.Context(), models.SeatsStruct{}, idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
		return
	}

	// if len(seats) < 1 {
	// 	ctx.JSON(http.StatusNotFound, gin.H{
	// 		"msg": "seats not found",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": seats,
	})
}