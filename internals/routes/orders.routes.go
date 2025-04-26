package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterOrders(router *gin.Engine, ordersRepo *repositories.OrdersRepository) {
	routerOrders := router.Group("/order/:id")
	ordersHandler := handlers.NewOrdersHandler(ordersRepo)

	// router create order
	routerOrders.POST("", ordersHandler.CreateOrder)
}