package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterOrders(router *gin.Engine, ordersRepo *repositories.OrdersRepository) {
	routerOrders := router.Group("/order")
	ordersHandler := handlers.NewOrdersHandler(ordersRepo)

	// router create order
	routerOrders.POST("/:id", ordersHandler.CreateOrder)

	// router get order history
	routerOrders.GET("/history/:id", ordersHandler.GetOrderHistory)
}