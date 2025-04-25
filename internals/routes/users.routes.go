package routes

import (
	"tikcitz-app/internals/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouterUsers(router *gin.Engine) {
	routerUsers := router.Group("/users")
	usersHandler := handlers.NewUsersHandlers()

	// add user
	routerUsers.POST("/signup", usersHandler.UserRegister)

	// auth user login
	routerUsers.POST("/login", usersHandler.UserLogin)

	// Get data profile by id
	routerUsers.GET("/profile/:id", usersHandler.GetProfileById)

	// Update data profile
	routerUsers.POST("/profile/:id", usersHandler.UpdateProfile)

}