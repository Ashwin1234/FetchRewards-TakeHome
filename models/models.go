package models

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
