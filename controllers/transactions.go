package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"fetch_rewards/models"

	"github.com/gin-gonic/gin"
)

var transactions []models.Transaction
var spend = make(map[string]int)
var balance = make(map[string]int)

/* struct initializations */

func UpdatePointBalance() {

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
func AddTransactions(c *gin.Context) {
	var transaction models.Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " Transaction request should be of specified format "})
		return
	}

	transactions = append(transactions, transaction)
	c.JSON(http.StatusCreated, gin.H{"data": transaction})

}

/* Service to get the points spent by each payer */

func GetSpendPoints(c *gin.Context) {

	UpdatePointBalance()

	var total = 0

	var spentList []models.PointsSpent

	var spendPoints models.SpendPoints

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
		u, err := json.Marshal(models.PointsSpent{Payer: key, Points: value})
		if err != nil {
			fmt.Println(err)
		}
		var jsonMap models.PointsSpent
		json.Unmarshal([]byte(u), &jsonMap)
		spentList = append(spentList, jsonMap)
	}

	c.JSON(http.StatusOK, spentList)

}

/* Service to get balance of each payer */

func GetPointBalances(c *gin.Context) {

	c.JSON(http.StatusOK, balance)
	return
}
