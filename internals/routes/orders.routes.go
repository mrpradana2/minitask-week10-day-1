package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/middleware"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterOrders(router *gin.Engine, ordersRepo *repositories.OrdersRepository, middleware *middleware.Middleware) {
	routerOrders := router.Group("/order")
	ordersHandler := handlers.NewOrdersHandler(ordersRepo)

	// router create order
	routerOrders.POST("", middleware.VerifyToken, middleware.AcceessGate("user"), ordersHandler.CreateOrder)

	// router get order history
	routerOrders.GET("", middleware.VerifyToken, middleware.AcceessGate("user"), ordersHandler.GetOrderHistory)

	// get order by id
}