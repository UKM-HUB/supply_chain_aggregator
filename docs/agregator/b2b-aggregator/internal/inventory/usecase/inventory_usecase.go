package usecase

import (
	"context"
	"b2b-aggregator/internal/inventory/repository"
)

type InventoryUsecase struct {
	geoRepo *repository.GeoRepository
	// Di industri nyata, tambahkan productRepo *repository.ProductRepository di sini
	// untuk mengecek jumlah stok asli di PostgreSQL
}

func NewInventoryUsecase(geoRepo *repository.GeoRepository) *InventoryUsecase {
	return &InventoryUsecase{geoRepo: geoRepo}
}

// Struct balasan internal
type SupplyResult struct {
	UMKMID   string
	Stock    int32
	Distance float64
}

func (u *InventoryUsecase) SearchNearbySupplies(ctx context.Context, productCode string, lat, lng, radiusKm float64) ([]SupplyResult, error) {
	// 1. Panggil Redis GEO melalui Repository
	locations, err := u.geoRepo.GetNearbyUMKM(ctx, lng, lat, radiusKm)
	if err != nil {
		return nil, err
	}

	var results []SupplyResult
	for _, loc := range locations {
		// 2. Gabungkan data
		// Idealnya di sini kita melakukan pengecekan ke DB: "Apakah UMKM ini punya stok productCode tersebut?"
		// Untuk MVP ini, kita asumsikan UMKM yang terdeteksi di Redis memiliki stok 100 (Mocking)
		results = append(results, SupplyResult{
			UMKMID:   loc.Name, // Name di Redis menyimpan UMKM_ID
			Stock:    100,      // Asumsi stok
			Distance: loc.Dist,
		})
	}

	return results, nil
}
