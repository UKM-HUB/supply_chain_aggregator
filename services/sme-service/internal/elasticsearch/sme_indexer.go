package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	pkges "supply-chain-aggregator/pkg/elasticsearch"
	"supply-chain-aggregator/services/sme-service/internal/entity"
)

const indexName = "smes"

// SMEDocument adalah representasi dokumen SME di Elasticsearch.
type SMEDocument struct {
	ID          string   `json:"id"`
	OwnerID     string   `json:"owner_id"`
	Name        string   `json:"name"`
	Phone       string   `json:"phone"`
	Address     string   `json:"address"`
	Description string   `json:"description"`
	CategoryIDs []string `json:"category_ids"`
	Products    []string `json:"products"`
	Capacity    string   `json:"capacity"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Status      string   `json:"status"`
}

// SMEIndexer mengelola indexing dan pencarian SME di Elasticsearch.
type SMEIndexer struct {
	client *pkges.Client
}

// NewSMEIndexer membuat indexer baru dan memastikan index sudah ada.
func NewSMEIndexer(client *pkges.Client) (*SMEIndexer, error) {
	indexer := &SMEIndexer{client: client}
	if err := indexer.ensureIndex(context.Background()); err != nil {
		return nil, err
	}
	return indexer, nil
}

// ensureIndex membuat index 'smes' dengan mapping yang sesuai jika belum ada.
func (i *SMEIndexer) ensureIndex(ctx context.Context) error {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":           map[string]interface{}{"type": "keyword"},
				"owner_id":     map[string]interface{}{"type": "keyword"},
				"name":         map[string]interface{}{"type": "text", "analyzer": "standard"},
				"phone":        map[string]interface{}{"type": "keyword"},
				"address":      map[string]interface{}{"type": "text"},
				"description":  map[string]interface{}{"type": "text", "analyzer": "standard"},
				"category_ids": map[string]interface{}{"type": "keyword"},
				"products":     map[string]interface{}{"type": "text"},
				"capacity":     map[string]interface{}{"type": "keyword"},
				"latitude":     map[string]interface{}{"type": "float"},
				"longitude":    map[string]interface{}{"type": "float"},
				"status":       map[string]interface{}{"type": "keyword"},
			},
		},
	}
	return i.client.CreateIndexIfNotExists(ctx, indexName, mapping)
}

// Index menyimpan atau memperbarui dokumen SME di Elasticsearch.
func (i *SMEIndexer) Index(ctx context.Context, sme entity.SME) error {
	doc := SMEDocument{
		ID:          sme.ID,
		OwnerID:     sme.OwnerID,
		Name:        sme.Name,
		Phone:       sme.Phone,
		Address:     sme.Address,
		Description: sme.Description,
		CategoryIDs: sme.CategoryIDs,
		Products:    sme.Products,
		Capacity:    sme.Capacity,
		Latitude:    sme.Latitude,
		Longitude:   sme.Longitude,
		Status:      sme.Status,
	}
	return i.client.IndexDocument(ctx, indexName, sme.ID, doc)
}

// Delete menghapus dokumen SME dari Elasticsearch berdasarkan ID.
func (i *SMEIndexer) Delete(ctx context.Context, smeID string) error {
	return i.client.DeleteDocument(ctx, indexName, smeID)
}

// SearchInput adalah parameter untuk full-text search SME.
type SearchInput struct {
	Query      string // full-text query terhadap name dan description
	CategoryID string
	Status     string
	From       int
	Size       int
}

// SearchResult adalah hasil pencarian SME dari Elasticsearch.
type SearchResult struct {
	Data  []SMEDocument
	Total int
}

// Search melakukan full-text search SME di Elasticsearch.
func (i *SMEIndexer) Search(ctx context.Context, input SearchInput) (*SearchResult, error) {
	if input.Size == 0 {
		input.Size = 20
	}

	must := make([]map[string]interface{}, 0)
	filter := make([]map[string]interface{}, 0)

	// Full-text search pada name dan description
	if input.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  input.Query,
				"fields": []string{"name^2", "description", "products"},
			},
		})
	}

	// Filter berdasarkan category_id
	if input.CategoryID != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"category_ids": input.CategoryID},
		})
	}

	// Filter berdasarkan status
	if input.Status != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"status": input.Status},
		})
	}

	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}
	if len(must) == 0 && len(filter) == 0 {
		boolQuery["must"] = []map[string]interface{}{{"match_all": map[string]interface{}{}}}
	}

	query := map[string]interface{}{
		"from":  input.From,
		"size":  input.Size,
		"query": map[string]interface{}{"bool": boolQuery},
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
		},
	}

	result, err := i.client.Search(ctx, indexName, query)
	if err != nil {
		return nil, fmt.Errorf("sme elasticsearch search: %w", err)
	}

	docs := make([]SMEDocument, 0, len(result.Hits))
	for _, hit := range result.Hits {
		var doc SMEDocument
		if err := json.Unmarshal(hit, &doc); err != nil {
			return nil, fmt.Errorf("sme elasticsearch unmarshal: %w", err)
		}
		docs = append(docs, doc)
	}

	return &SearchResult{
		Data:  docs,
		Total: result.Total,
	}, nil
}
