package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type UMKM struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Distance float64 `json:"distance_in_meter"`
}

var db *pgxpool.Pool

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to connect database:", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Database connected")
}

func nearbyUMKM(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat and lng are required",
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid latitude",
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid longitude",
		})
		return
	}

	query := `
	SELECT 
		id,
		name,
		address,
		ST_Distance(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		) AS distance
	FROM umkms
	ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
	LIMIT 10;
	`

	rows, err := db.Query(context.Background(), query, lng, lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var umkms []UMKM

	for rows.Next() {
		var u UMKM

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Address,
			&u.Distance,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		umkms = append(umkms, u)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    umkms,
	})
}

func main() {
	initDB()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/nearby/umkm", nearbyUMKM)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}