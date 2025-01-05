package routes

import(
	controller "github.com/rahulgit24/golang-projects/auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup",controller.SignUp())
	incomingRoutes.POST("users/login",controller.Login())
}	