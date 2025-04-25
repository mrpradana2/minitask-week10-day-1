package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterUsers(router *gin.Engine, usersRepo *repositories.UserRepository) {
	routerUsers := router.Group("/users")
	usersHandler := handlers.NewUsersHandlers(usersRepo)

	// router add user
	routerUsers.POST("/signup", usersHandler.UserRegister)

	// router auth user login
	routerUsers.POST("/login", usersHandler.UserLogin)

	// router Get data profile by id
	routerUsers.GET("/profile/:id", usersHandler.GetProfileById)

	// router Update data profile
	routerUsers.PUT("/profile/:id", usersHandler.UpdateProfile)

}