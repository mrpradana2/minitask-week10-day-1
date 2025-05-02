package routes

import (
	"tikcitz-app/internals/handlers"
	"tikcitz-app/internals/middleware"
	"tikcitz-app/internals/repositories"

	"github.com/gin-gonic/gin"
)

func InitRouterUsers(router *gin.Engine, usersRepo *repositories.UserRepository, middleware *middleware.Middleware) {
	routerUsers := router.Group("/users")
	usersHandler := handlers.NewUsersHandlers(usersRepo)



	// router add user (fix)
	routerUsers.POST("/signup", usersHandler.UserRegister)

	// router auth user login (fix)
	routerUsers.POST("/login", usersHandler.UserLogin)

	// router Get data profile by id (fix)
	routerUsers.GET("/profile/:id", usersHandler.GetProfileById)

	// router Update data profile (fix)
	routerUsers.PATCH("/profile/", middleware.VerifyToken, middleware.AcceessGate("user"), usersHandler.UpdateProfile)

	// // router untuk verify user token
	// routerUsers.GET("/verify", usersHandler.VerifyToken)

}