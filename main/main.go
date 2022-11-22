package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

var transactions []Transaction
var spend = make(map[string]int)
var balance = make(map[string]int)

func main() {

	r := gin.New()
	r.GET("/spend_points", getSpendPoints)
	r.POST("/add_transactions", addTransactions)
	r.GET("/point_balances", getPointBalances)

	r.Run()

}

type Transaction struct {
	Payer     string `json:"payer" binding:"required"`
	Points    int    `json:"points" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}

type SpendPoints struct {
	Points int `json:"points" binding:"required"`
}

type PointsSpent struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}

func addTransactions(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " Transaction request should be of specified format "})
		return
	}

	transactions = append(transactions, transaction)
	if val, ok := balance[transaction.Payer]; ok {
		balance[transaction.Payer] = val + transaction.Points
	} else {
		balance[transaction.Payer] = transaction.Points
	}
	c.JSON(http.StatusOK, gin.H{"data": transaction})

}

func getSpendPoints(c *gin.Context) {

	/*(var balance = make(map[string]int)

	var spendMap = make(map[string]int)

	var tempTransactions = transactions

	var spendPoints SpendPoints

	var pointsSpent []string

	if err := c.BindJSON(&spendPoints); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var points_to_spend int = spendPoints.Points

	sort.Slice(tempTransactions, func(i, j int) bool {
		return tempTransactions[i].Timestamp < tempTransactions[j].Timestamp
	})

	for _, value := range tempTransactions {
		if value.Points > 0 {
			var pointsForPayer = 0
			for _, value1 := range tempTransactions {
				if value1.Payer == value.Payer {
					if value1.Points < 0 {
						pointsForPayer = pointsForPayer + value.Points - value1.Points
					}
				}
			}
			if pointsForPayer > 0 {
				if points_to_spend-pointsForPayer > 0 {
					spendMap[value.Payer] = spendMap[value.Payer] - pointsForPayer
					points_to_spend = points_to_spend - pointsForPayer
				} else {
					spendMap[value.Payer] = spendMap[value.Payer] - points_to_spend
					points_to_spend = 0
				}

			}
		}
	}

	for key, _ := range spendMap {
		balance[key] = balance[key] - spendMap[key]
	}

	for key, value := range spendMap {
		payerValue := &PointsSpent{
			Key:   key,
			Value: value,
		}
		data, _ := json.Marshal(payerValue)
		pointsSpent = append(pointsSpent, string(data))
	}

	spend = balance */

	var total = 0

	var spentList []string

	var spendPoints SpendPoints

	if err := c.BindJSON(&spendPoints); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spend point request should be of specified format"})
		return
	}

	var spend = spendPoints.Points

	var spendMap = make(map[string]int)

	if spend < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spend point cannot be negative"})
		return
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp < transactions[j].Timestamp
	})

	for _, transaction := range transactions {

		if total < spend {
			if total+transaction.Points <= spend {
				total = total + transaction.Points

				if val, ok := spendMap[transaction.Payer]; ok {
					spendMap[transaction.Payer] = val - transaction.Points
				} else {
					spendMap[transaction.Payer] = -transaction.Points
				}
				balance[transaction.Payer] = balance[transaction.Payer] - transaction.Points

			} else {

				if val, ok := spendMap[transaction.Payer]; ok {
					spendMap[transaction.Payer] = val - (spend - total)
				} else {
					spendMap[transaction.Payer] = -(spend - total)
				}
				balance[transaction.Payer] = balance[transaction.Payer] - (spend - total)
				break
			}
		} else {
			break
		}
	}

	for key, value := range spendMap {
		u, err := json.Marshal(PointsSpent{Payer: key, Points: value})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(u))
		spentList = append(spentList, string(u))
	}

	c.JSON(http.StatusOK, spendMap)

}

func getPointBalances(c *gin.Context) {

	c.JSON(http.StatusOK, balance)
}
