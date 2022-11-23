package routes

import (
	"fetch_rewards/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.New()
	r.GET("/spend_points", controllers.GetSpendPoints)
	r.POST("/add_transactions", controllers.AddTransactions)
	r.GET("/point_balances", controllers.GetPointBalances)

	return r
}
