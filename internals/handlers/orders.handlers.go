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

	errCreateOrder := o.ordersRepo.CreateOrder(ctx.Request.Context(), newOrder, idInt)

	if errCreateOrder != nil {
		log.Println("Insert profile error:", errCreateOrder)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	// if cmd.RowsAffected() == 0 {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "add data successfully", 
	})
}

// Handler get order history user
func (o *OrdersHandler) GetOrderHistory(ctx *gin.Context) {
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

	result, err := o.ordersRepo.GetOrderHistory(ctx.Request.Context(), idInt)

	if err != nil {
		log.Println("[ERROR]: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occured on the server",
		})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "order not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
	})
}