package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	de "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	"github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/geocodeerrors"
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

func (p *Processor) Process(ctx context.Context, rawEvent models.RawOsintEvent) error {
	var processedEvent storage.ProcessedEvent
	var err error

	switch event := rawEvent.(type) {
	case models.RawTelegramEvent:
		processedEvent, err = p.processTelegramEvent(event, ctx)
	default:
		err = fmt.Errorf("Unsupported event type given: %T", rawEvent)
	}

	err, docId := p.ESClient.IndexEvent(ctx, processedEvent)
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

func (p *Processor) handleGeocodeError(err *geocodeerrors.GeocodeError) error {
	var result error = nil

	switch err.Code {
	case geocodeerrors.ErrCodeNetworkError:
	case geocodeerrors.ErrCodeInvalidRequest:
	case geocodeerrors.ErrCodeInternalError:
	case geocodeerrors.ErrCodeParsingError:
		result = fmt.Errorf("Error of type %s occoured while geocoding: %v", err.Code, err)
	case geocodeerrors.ErrCodeFiltered:
	case geocodeerrors.ErrCodeNotFound:
		// If a message is filtered / not found then the processor should not
		// return an error - its the inteded behaviour.
		result = nil
	}

	return result
}

func (p *Processor) processTelegramEvent(te models.RawTelegramEvent, ctx context.Context) (storage.ProcessedEvent, error) {
	var processedEvent storage.ProcessedEvent

	result, err := de.CreateAgentSummary(te.Text, ctx, p.Cfg.OpenAIKey, p.Cfg.OpenAIModel)
	if err != nil {
		return processedEvent, fmt.Errorf("error extracting info: %w", err)
	}

	log.Printf("AI Summary: %+v\n", result)

	locationMap, geocodeErr := p.Geocoder.GetBatchCoordinates(result.EnLocations)

	if geocodeErr != nil {
		return processedEvent, p.handleGeocodeError(geocodeErr)
	}

	var locations []models.Location
	for name, geo := range locationMap {
		locations = append(locations, models.Location{
			Name: name,
			Lat:  geo.Lat,
			Lon:  geo.Lon,
		})
	}

	processedEvent = storage.ProcessedEvent{
		RawMessage:     te.Text,
		Summary:        result.HeSummary,
		Locations:      locations,
		TimestampEpoch: models.ParseToEpoch(te.Timestamp),
	}

	return processedEvent, nil
}
