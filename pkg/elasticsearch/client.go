package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps the official Elasticsearch client.
type Client struct {
	es *elasticsearch.Client
}

// Config holds Elasticsearch connection settings.
type Config struct {
	Addresses []string // e.g. ["http://localhost:9200"]
	Username  string
	Password  string
}

// New creates a new Elasticsearch client and verifies connectivity.
func New(cfg Config) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: failed to create client: %w", err)
	}

	// Verify connection
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: failed to connect: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch: info response error: %s", res.Status())
	}

	return &Client{es: es}, nil
}

// IndexDocument indexes a single document into the given index.
// If docID is empty, Elasticsearch will generate one.
func (c *Client) IndexDocument(ctx context.Context, index, docID string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal error: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch: index error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch: index response [%s]: %s", res.Status(), string(body))
	}
	return nil
}

// DeleteDocument removes a document from an index by its ID.
func (c *Client) DeleteDocument(ctx context.Context, index, docID string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("elasticsearch: delete error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch: delete response [%s]: %s", res.Status(), string(body))
	}
	return nil
}

// SearchResult is returned by Search.
type SearchResult struct {
	Total int
	Hits  []json.RawMessage
}

// Search executes a query DSL against the given index.
// query should be a valid Elasticsearch Query DSL map.
func (c *Client) Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResult, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: marshal query error: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(strings.NewReader(string(body))),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		rawBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch: search response [%s]: %s", res.Status(), string(rawBody))
	}

	var esResp struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, fmt.Errorf("elasticsearch: decode response error: %w", err)
	}

	hits := make([]json.RawMessage, 0, len(esResp.Hits.Hits))
	for _, h := range esResp.Hits.Hits {
		hits = append(hits, h.Source)
	}

	return &SearchResult{
		Total: esResp.Hits.Total.Value,
		Hits:  hits,
	}, nil
}

// CreateIndexIfNotExists creates an index with the given mapping only if it doesn't already exist.
func (c *Client) CreateIndexIfNotExists(ctx context.Context, index string, mapping map[string]interface{}) error {
	existsRes, err := c.es.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("elasticsearch: check index exists error: %w", err)
	}
	defer existsRes.Body.Close()

	if existsRes.StatusCode == 200 {
		// Index already exists
		return nil
	}

	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal mapping error: %w", err)
	}

	createRes, err := c.es.Indices.Create(index,
		c.es.Indices.Create.WithBody(bytes.NewReader(body)),
		c.es.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch: create index error: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		rawBody, _ := io.ReadAll(createRes.Body)
		return fmt.Errorf("elasticsearch: create index response [%s]: %s", createRes.Status(), string(rawBody))
	}
	return nil
}
