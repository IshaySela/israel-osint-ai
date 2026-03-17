package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dataextraction "github.com/IshaySela/israel-osint-ai/services/processing/data-extraction"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

type ProcessedEvent struct {
	RawMessage string                            `json:"raw_message"`
	Summary    string                            `json:"summary"`
	Locations  map[string]dataextraction.Geocode `json:"locations"`
	Timestamp  string                            `json:"timestamp"`
}

func NewElasticsearchClient() *ElasticsearchClient {
	return &ElasticsearchClient{}
}

func (esc *ElasticsearchClient) Setup(addresses []string) error {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating the elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("error getting elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from elasticsearch: %s", res.String())
	}

	esc.client = client
	return nil
}

func (esc *ElasticsearchClient) IndexEvent(ctx context.Context, index string, event ProcessedEvent) error {
	if esc.client == nil {
		return fmt.Errorf("elasticsearch client not initialized, call Setup first")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	res, err := esc.client.Index(
		index,
		bytes.NewReader(data),
		esc.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document in elasticsearch: %s", res.String())
	}

	return nil
}
