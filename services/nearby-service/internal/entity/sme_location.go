package entity

type SMELocation struct {
	ID          string
	Name        string
	Address     string
	Description string
	CategoryIDs []string
	Latitude    float64
	Longitude   float64
	Status      string
}

type NearbySME struct {
	SMELocation
	DistanceKM float64
}
