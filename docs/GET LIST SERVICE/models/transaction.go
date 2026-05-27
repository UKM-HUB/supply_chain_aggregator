package models

type Transaction struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Invoice  string `json:"invoice"`
	Amount   int    `json:"amount"`
	Status   string `json:"status"`
}