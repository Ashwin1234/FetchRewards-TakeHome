package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"fetch_rewards/models"

	"github.com/gin-gonic/gin"
)

var input []models.Transaction
var spend = make(map[string]int)
var balance = make(map[string]int)
var transactions []models.Transaction

/* func to update transaction based on timestamp */
func UpdateTransaction() {
	transactions = nil
	sort.Slice(input, func(i, j int) bool {
		return input[i].Timestamp < input[j].Timestamp
	})

	for i := 0; i < len(input); i++ {
		if input[i].Points > 0 {
			transactions = append(transactions, input[i])
		}
	}

	for i := len(input) - 1; i >= 0; i-- {
		if input[i].Points < 0 {
			transactions = append([]models.Transaction{input[i]}, transactions...)
		}
	}

}

/* func to update point balance */

func UpdatePointBalance() {

	balance = make(map[string]int)

	for _, transaction := range input {
		if val, ok := balance[transaction.Payer]; ok {
			balance[transaction.Payer] = val + transaction.Points
		} else {
			balance[transaction.Payer] = transaction.Points
		}
	}

}

/* function to calculate the total points of all payers. */

func getTotalPoints() int {
	var total = 0
	for _, value := range balance {
		total = total + value
	}

	return total
}

/* Service to add transactions */
func AddTransactions(c *gin.Context) {
	var transaction models.Transaction
	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " Transaction request should be of specified format "})
		return
	}

	input = append(input, transaction)
	c.JSON(http.StatusCreated, gin.H{"success": "transaction added"})

	UpdatePointBalance()

}

/* Service to get the points spent by each payer */

func GetSpendPoints(c *gin.Context) {

	var remainingPoints = getTotalPoints()

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spend points cannot be negative"})
		return
	}

	if spend > remainingPoints {
		c.JSON(http.StatusBadRequest, gin.H{"error": " Spend points cannot be more than balance"})
		return
	}

	UpdateTransaction()

	for _, transaction := range transactions {

		if total < spend {

			if total+transaction.Points <= spend || total+balance[transaction.Payer] <= spend {

				if balance[transaction.Payer] < transaction.Points {

					total = total + balance[transaction.Payer]
					if val, ok := spendMap[transaction.Payer]; ok {
						spendMap[transaction.Payer] = val - balance[transaction.Payer]
					} else {
						spendMap[transaction.Payer] = -balance[transaction.Payer]
					}
					balance[transaction.Payer] = 0

				} else {
					total = total + transaction.Points

					if val, ok := spendMap[transaction.Payer]; ok {
						spendMap[transaction.Payer] = val - transaction.Points
					} else {
						spendMap[transaction.Payer] = -transaction.Points
					}
					balance[transaction.Payer] = balance[transaction.Payer] - transaction.Points
				}

			} else {

				if balance[transaction.Payer] < (spend - total) {

					total = balance[transaction.Payer]
					if val, ok := spendMap[transaction.Payer]; ok {
						spendMap[transaction.Payer] = val - balance[transaction.Payer]
					} else {
						spendMap[transaction.Payer] = -balance[transaction.Payer]
					}
					balance[transaction.Payer] = 0

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
			}

		} else {
			break
		}
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
