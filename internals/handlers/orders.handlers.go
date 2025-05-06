package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	ordersRepo *repositories.OrdersRepository
}

func NewOrdersHandler(ordersRepo *repositories.OrdersRepository) *OrdersHandler {
	return &OrdersHandler{ordersRepo: ordersRepo}
}

// Handler create order (fix)
func (o *OrdersHandler) CreateOrder(ctx *gin.Context) {
	// ambil id yang ada di header
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// buat variable untuk menampung data movie baru dari admin
	var newOrder models.OrdersStr

	// binding data
	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid data sent",
		})
		return
	}

	errCreateOrder := o.ordersRepo.CreateOrder(ctx.Request.Context(), newOrder, userClaims.Id)

	if errCreateOrder != nil {
		log.Println("[DEBUG] : ", errCreateOrder)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "succes create order",
	})
}

// Handler get order history user (fix)
func (o *OrdersHandler) GetOrderHistory(ctx *gin.Context) {

	// mengambil data berdasarkan user yang login 
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)
	idInt := userClaims.Id

	// menjalankan fungsi repository get order history
	result, err := o.ordersRepo.GetOrderHistory(ctx.Request.Context(), idInt)

	// error handling jika gagal menjalankan query
	if err != nil {
		log.Println("[ERROR]: ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "server error",
		})
		return
	}

	// mengecek jika order 0, maka tampilkan error movie not found
	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "failed",
			Msg: "orders not found",
		})
		return
	}

	// jika berhasil mengambil data dari server dan ada datanya, maka tampilkan pesan ini
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
}

// Handler get order by order_id 
func (o *OrdersHandler) GetOrderById(ctx *gin.Context) {

	// mengambil data berdasarkan user yang login 
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// ambil id param
	idStr, ok := ctx.Params.Get("orderId")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "Param id is needed",
		})
		return
	}

	// konversi id string menjadi id integer
	orderId, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[BEDUG] : ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	result, err := o.ordersRepo.GetOrderById(ctx.Request.Context(), userClaims.Id, orderId)
	// error handling jika gagal menjalankan query
	if err != nil {
		log.Println("[ERROR]: ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "server error",
		})
		return
	}

	// mengecek jika order 0, maka tampilkan error movie not found
	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "failed",
			Msg: "orders not found",
		})
		return
	}

	// jika berhasil mengambil data dari server dan ada datanya, maka tampilkan pesan ini
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
} 