package models

type Factory struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Status string `json:"status"`
}