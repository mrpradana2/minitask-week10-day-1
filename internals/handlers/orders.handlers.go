package handlers

import (
	"errors"
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

// Create Order
// @summary					Create order movie
// @router					/order [post]
// @Description 			Create order seat movie
// @Tags        			Order
// @Param        			requestBody body models.RequestOrdersStruct true "Input data for create order"
// @Param 					Authorization header string true "Bearer Token"
// @accept					json
// @produce					json
// @Security    			BearerAuth
// @failure					500 {object} models.MessageInternalServerError
// @failure					400 {object} models.MessageBadRequest
// @failure					409 {object} models.MessageConflict
// @success					201 {object} models.MessageCreated
func (o *OrdersHandler) CreateOrder(ctx *gin.Context) {
	// ambil id yang ada di header
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// buat variable untuk menampung data movie baru dari admin
	var newOrder models.OrdersStr

	// binding data
	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "invalid data sent",
		})
		return
	}

	errCreateOrder := o.ordersRepo.CreateOrder(ctx.Request.Context(), newOrder, userClaims.Id)

	if errCreateOrder != nil {
		log.Println("[ERROR] : ", errCreateOrder.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Message{
		Status: http.StatusCreated,
		Msg: "succes create order",
	})
}

// Get Order user
// @summary					Get order history
// @router					/order [get]
// @Description 			Get history order user 
// @Tags        			Order
// @Param 					Authorization header string true "Bearer Token"
// @Accept					json
// @produce					json
// @Security    			BearerAuth
// @failure					500 {object} models.MessageInternalServerError
// @failure					404 {object} models.MessageNotFound
// @failure					409 {object} models.MessageConflict
// @success					200 {object} models.MessageResult
func (o *OrdersHandler) GetOrderHistory(ctx *gin.Context) {
	
	// mengambil data berdasarkan user yang login 
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)
	idInt := userClaims.Id
	
	// menjalankan fungsi repository get order history
	result, err := o.ordersRepo.GetOrderHistory(ctx.Request.Context(), idInt)
	
	// error handling jika gagal menjalankan query
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "server error",
		})
		return
	}
	
	// mengecek jika order 0, maka tampilkan error movie not found
	if len(result) < 1 {
		log.Println("[ERROR] : ", errors.New("orders not found"))
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "orders not found",
		})
		return
	}

	// jika berhasil mengambil data dari server dan ada datanya, maka tampilkan pesan ini
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "success",
		Result: result,
	})
}


// Get Order History by id
// @summary					Get order history
// @router					/order/:orderId [get]
// @Description 			Get history order user by id order  
// @Tags        			Order
// @Param        			orderId query string true "order id for get detail history order"
// @Param 					Authorization header string true "Bearer Token"
// @Accept					json
// @produce					json
// @Security    			BearerAuth
// @failure					500 {object} models.MessageInternalServerError 
// @failure					404 {object} models.MessageNotFound
// @failure					409 {object} models.MessageConflict
// @success					200 {object} models.MessageResult
func (o *OrdersHandler) GetOrderById(ctx *gin.Context) {

	// mengambil data berdasarkan user yang login 
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	// ambil id param
	idStr, ok := ctx.Params.Get("orderId")

	// handling error jika param tidak ada
	if !ok {
		log.Println("[ERROR] : ", errors.New("params not found"))
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: http.StatusBadRequest,
			Msg: "Param id is needed",
		})
		return
	}

	// konversi id string menjadi id integer
	orderId, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "an error occurred on the server",
		})
		return
	}

	result, err := o.ordersRepo.GetOrderById(ctx.Request.Context(), userClaims.Id, orderId)
	// error handling jika gagal menjalankan query
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: http.StatusInternalServerError,
			Msg: "server error",
		})
		return
	}

	// mengecek jika order 0, maka tampilkan error movie not found
	if len(result) < 1 {
		log.Println("[ERROR] : ", errors.New("orders not found"))
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: http.StatusNotFound,
			Msg: "orders not found",
		})
		return
	}

	// jika berhasil mengambil data dari server dan ada datanya, maka tampilkan pesan ini
	ctx.JSON(http.StatusOK, models.Message{
		Status: http.StatusOK,
		Msg: "success",
		Result: result,
	})
} 