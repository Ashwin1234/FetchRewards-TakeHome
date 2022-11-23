package controllers

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

/* struct initializations */

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

func updatePointBalance() {

	balance = make(map[string]int)

	for _, transaction := range transactions {
		if val, ok := balance[transaction.Payer]; ok {
			balance[transaction.Payer] = val + transaction.Points
		} else {
			balance[transaction.Payer] = transaction.Points
		}
	}

	fmt.Println(balance)
}

/* Service to add transactions */
func addTransactions(c *gin.Context) {
	var transaction Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " Transaction request should be of specified format "})
		return
	}

	transactions = append(transactions, transaction)
	c.JSON(http.StatusCreated, gin.H{"data": transaction})

}

/* Service to get the points spent by each payer */

func getSpendPoints(c *gin.Context) {

	updatePointBalance()

	var total = 0

	var spentList []PointsSpent

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
				total = spend
				break
			}
		} else {
			break
		}
	}

	if total < spend {
		c.JSON(http.StatusOK, gin.H{"error": "Insufficient balance"})
		return
	}

	for key, value := range spendMap {
		u, err := json.Marshal(PointsSpent{Payer: key, Points: value})
		if err != nil {
			fmt.Println(err)
		}
		var jsonMap PointsSpent
		json.Unmarshal([]byte(u), &jsonMap)
		spentList = append(spentList, jsonMap)
	}

	c.JSON(http.StatusOK, spentList)

}
