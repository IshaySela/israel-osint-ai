package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ElasticsearchClient struct {
	client *elasticsearch.TypedClient
	cfg    *config.Config
}

type ProcessedEvent struct {
	RawMessage     string            `json:"raw_message"`
	Summary        string            `json:"summary"`
	Locations      []models.Location `json:"locations"`
	TimestampEpoch int64             `json:"timestamp_epoch"`
	ChannelTitle   string            `json:"channel_title"`
}

func NewElasticsearchClient(cfg *config.Config) *ElasticsearchClient {
	return &ElasticsearchClient{
		cfg: cfg,
	}
}

func (esc *ElasticsearchClient) Setup(addresses []string) error {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating the elasticsearch client: %w", err)
	}
	_, err = client.Ping().Do(context.TODO())

	if err != nil {
		return fmt.Errorf("Error connecting to elasticsearch: %w", err)
	}

	esc.client = client
	return nil
}

func (esc *ElasticsearchClient) IndexEvent(ctx context.Context, event ProcessedEvent) (error, string) {
	if esc.client == nil {
		return fmt.Errorf("elasticsearch client not initialized, call Setup first"), ""
	}

	res, err := esc.client.
		Index(esc.cfg.ProcessedEventsIndex).
		Document(event).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error indexing event to elasticsearch: %w", err), ""
	}

	return nil, res.Id_
}

func (esc *ElasticsearchClient) IndexGeocode(ctx context.Context, locationText string, geocode models.Geocode) (error, string) {
	if esc.client == nil {
		return fmt.Errorf("elasticsearch client not initialized, call Setup first"), ""
	}

	lat, _ := strconv.ParseFloat(geocode.Lat, 64)
	lon, _ := strconv.ParseFloat(geocode.Lon, 64)

	cache := models.GeocodeCache{
		LocationText: locationText,
		Lat:          lat,
		Lon:          lon,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	res, err := esc.client.
		Index(esc.cfg.ElasticsearchGeocodeIndex).
		Document(cache).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error indexing geocode cache: %w", err), ""
	}

	return nil, res.Id_
}

func (esc *ElasticsearchClient) GetGeocode(ctx context.Context, location string) (models.GeocodeCache, error) {
	if esc.client == nil {
		return models.GeocodeCache{}, fmt.Errorf("elasticsearch client not initialized, call Setup first")
	}

	searchResult, err := esc.client.
		Search().
		Index(esc.cfg.ElasticsearchGeocodeIndex).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"location_text": {Query: location},
				},
			},
		}).Do(ctx)

	if err != nil {
		return models.GeocodeCache{}, fmt.Errorf("Could not find location in cache")
	}

	if searchResult.Hits.Total.Value == 0 {
		return models.GeocodeCache{}, fmt.Errorf("Cache miss while retriving geocode from es")
	}
	var parsed models.GeocodeCache

	err = json.Unmarshal(searchResult.Hits.Hits[0].Source_, &parsed)
	if err != nil {
		return models.GeocodeCache{}, fmt.Errorf("Error while parsing result from es %s", err.Error())
	}

	return parsed, nil
}
