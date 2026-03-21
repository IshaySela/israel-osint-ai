package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	de "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

type ProcessedEvent struct {
	RawMessage string                `json:"raw_message"`
	Summary    string                `json:"summary"`
	Locations  map[string]de.Geocode `json:"locations"`
	Timestamp  string                `json:"timestamp"`
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

func (esc *ElasticsearchClient) IndexGeocode(ctx context.Context, index string, locationText string, geocode de.Geocode) error {
	if esc.client == nil {
		return fmt.Errorf("elasticsearch client not initialized, call Setup first")
	}

	lat, _ := strconv.ParseFloat(geocode.Lat, 64)
	lon, _ := strconv.ParseFloat(geocode.Lon, 64)

	cache := models.GeocodeCache{
		LocationText: locationText,
		Lat:          lat,
		Lon:          lon,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("error marshaling geocode cache: %w", err)
	}

	res, err := esc.client.Index(
		index,
		bytes.NewReader(data),
		esc.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error indexing geocode cache: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing geocode cache in elasticsearch: %s", res.String())
	}

	return nil
}
