package routes

import (
	"fetch_rewards/controllers"

	"github.com/gin-gonic/gin"
)

// set up routes
func SetUpRouter() *gin.Engine {
	r := gin.New()

	// route to spend points based on the rules provided given the points.
	r.GET("/spend_points", controllers.GetSpendPoints)

	// route to add transactions.
	r.POST("/add_transactions", controllers.AddTransactions)

	// route to point balances of all the payers.
	r.GET("/point_balances", controllers.GetPointBalances)

	return r
}
