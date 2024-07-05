package routers

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/topup", controllers.TopUp)
	r.POST("/pay", controllers.Payment)
	r.POST("/transfer", controllers.Transfer)
	r.GET("/transactions", controllers.Transactions)
	r.PUT("/profile", controllers.UpdateProfile)

	return r
}
