package routes

import (
	controller "github.com/rahulgit24/golang-projects/auth/controllers"
	"github.com/rahulgit24/golang-projects/auth/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	// incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/users/:userid",controller.GetUser())
}