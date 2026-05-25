package main

import (
	"log"
	"b2b-aggregator/internal/auth/entity"
	"b2b-aggregator/pkg/config"
)

func main() {
	db, err := config.InitDatabase("localhost", "root", "secretpassword", "b2b_db", "5432")
	if err != nil {
		log.Fatalf("Gagal connect DB: %v", err)
	}

	// Auto Migrate Table
	db.AutoMigrate(&entity.User{})

	// Insert Dummy Data
	mockUMKM := entity.User{
		Username:  "Toko Sembako Berkah",
		Email:     "berkah@umkm.com",
		Password:  "hashed_password", // Harusnya di bcrypt dulu
		Role:      "UMKM",
		Latitude:  -6.200000,
		Longitude: 106.816666,
	}

	if err := db.Create(&mockUMKM).Error; err != nil {
		log.Printf("Data mungkin sudah ada: %v", err)
	} else {
		log.Println("✅ Berhasil menyuntikkan data UMKM dummy ke Database!")
	}
}
