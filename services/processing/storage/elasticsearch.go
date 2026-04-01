package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	client *elasticsearch.TypedClient
}

type ProcessedEvent struct {
	RawMessage string                    `json:"raw_message"`
	Summary    string                    `json:"summary"`
	Locations  map[string]models.Geocode `json:"locations"`
	Timestamp  string                    `json:"timestamp"`
}

func NewElasticsearchClient() *ElasticsearchClient {
	return &ElasticsearchClient{}
}

func (esc *ElasticsearchClient) Setup(addresses []string) error {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating the elasticsearch client: %w", err)
	}
	_, err = esc.client.Ping().Do(context.TODO())

	if err != nil {
		return fmt.Errorf("Error connecting to elasticsearch: %w", err)
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

	_, err = esc.client.
		Index(index).
		Document(data).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error indexing event to elasticsearch: %w", err)
	}

	return nil
}

func (esc *ElasticsearchClient) IndexGeocode(ctx context.Context, index string, locationText string, geocode models.Geocode) error {
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

	_, err := esc.client.Index(index).Document(cache).Do(ctx)

	if err != nil {
		return fmt.Errorf("error indexing geocode cache: %w", err)
	}

	return nil
}
