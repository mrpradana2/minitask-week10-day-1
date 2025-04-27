package handlers

import (
	"net/http"
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
	cinemaQ := ctx.Query("cinema")

	if cinemaQ == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "query cinema is needed",
		})
		return
	}

	seats, err := s.seatsRepo.GetSeatsAvailable(ctx.Request.Context(), models.MoviesStruct{}, cinemaQ)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
	}

	if len(seats) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "seats not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": seats,
	})
}