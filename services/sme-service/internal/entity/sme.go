package entity

import "time"

type SME struct {
	ID             string
	OwnerID        string
	Name           string
	Phone          string
	Address        string
	Description    string
	CategoryIDs    []string
	Products       []string
	Capacity       string
	Latitude       float64
	Longitude      float64
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
