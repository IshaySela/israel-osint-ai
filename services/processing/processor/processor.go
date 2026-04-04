package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	de "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	mb "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

type Processor struct {
	Cfg      *config.Config
	Geocoder *de.GeocodingService
	ESClient *storage.ElasticsearchClient
	broker   *mb.RabbitClient
}

func NewProcessor(cfg *config.Config, geocoder *de.GeocodingService, esClient *storage.ElasticsearchClient, broker *mb.RabbitClient) *Processor {
	return &Processor{
		Cfg:      cfg,
		Geocoder: geocoder,
		ESClient: esClient,
		broker:   broker,
	}
}

func (p *Processor) Process(ctx context.Context, event models.RawOsintEvent) error {
	log.Printf("Processing event: %s\n", string(event.Text))
	result, err := de.CreateAgentSummary(event, ctx, p.Cfg.OpenAIKey, p.Cfg.OpenAIModel)
	if err != nil {
		return fmt.Errorf("error extracting info: %w", err)
	}

	log.Printf("AI Summary: %+v\n", result)

	locationMap, geocodeErr := p.Geocoder.GetBatchCoordinates(result.EnLocations)
	if geocodeErr != nil {
		return fmt.Errorf("error fetching coordinates: %w", geocodeErr)
	}

	for loc, geo := range locationMap {
		err, _ := p.ESClient.IndexGeocode(ctx, p.Cfg.ElasticsearchGeocodeIndex, loc, geo)
		if err != nil {
			log.Printf("Error indexing geocode for %s: %v\n", loc, err)
		}
	}

	var locations []models.Location
	for name, geo := range locationMap {
		locations = append(locations, models.Location{
			Name: name,
			Lat:  geo.Lat,
			Lon:  geo.Lon,
		})
	}

	processedEvent := storage.ProcessedEvent{
		RawMessage: event.Text,
		Summary:    result.HeSummary,
		Locations:  locations,
		Timestamp:  event.Date,
	}

	err, docId := p.ESClient.IndexEvent(ctx, p.Cfg.ElasticsearchIndex, processedEvent)
	if err != nil {
		return fmt.Errorf("error indexing event to elasticsearch: %w", err)
	}

	log.Println("Successfully indexed event to elasticsearch")

	if err = p.broker.PublishProcessedEvent(processedEvent, docId); err != nil {
		return fmt.Errorf("error publishing processed event: %w", err)
	}

	log.Printf("Published the processed event to %s", p.Cfg.ProcessedEventsExchange)
	return nil
}
