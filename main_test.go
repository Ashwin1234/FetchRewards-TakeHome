package main

import (
	"bytes"
	"fetch_rewards/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestAddTransaction(t *testing.T) {

	r := SetUpRouter()
	r.POST("/add_transactions", controllers.AddTransactions)
	transaction := controllers.Transaction{
		Payer:     "DANNON",
		Points:    1000,
		Timestamp: "2022-11-02T14:00:00Z",
	}

	jsonValue, _ := json.Marshal(transaction)
	req, _ := http.NewRequest("POST", "/add_transactions", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetSpendPoints(t *testing.T) {

	r := SetUpRouter()
	r.GET("/spend_points", controllers.GetSpendPoints)

	spendPoints := controllers.SpendPoints{
		Points: 5000,
	}

	jsonValue, _ := json.Marshal(spendPoints)
	req, _ := http.NewRequest("GET", "/spend_points", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, responseData)

}

func TestGetPointBalance(t *testing.T) {

	r := SetUpRouter()
	r.GET("/point_balances", controllers.GetPointBalances)

	req, _ := http.NewRequest("GET", "/point_balances", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, responseData)
}
