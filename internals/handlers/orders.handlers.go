package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	ordersRepo *repositories.OrdersRepository
}

func NewOrdersHandler(ordersRepo *repositories.OrdersRepository) *OrdersHandler {
	return &OrdersHandler{ordersRepo: ordersRepo}
}

// Handler create order
func (o *OrdersHandler) CreateOrder(ctx *gin.Context) {
	var newOrder models.OrdersStruct

	// binding data
	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Param id is needed",
		})
		return
	}

	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
		return
	}

	cmd, err := o.ordersRepo.CreateOrder(ctx.Request.Context(), newOrder, idInt)

	if err != nil {
		log.Println("Insert profile error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if cmd.RowsAffected() == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "add data successfully", 
	})
}