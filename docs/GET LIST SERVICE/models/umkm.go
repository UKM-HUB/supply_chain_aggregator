package models

type UMKM struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Status string `json:"status"`
}