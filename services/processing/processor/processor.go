package processor

import (
	"context"
	"log"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	dataextraction "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

type Processor struct {
	Cfg      *config.Config
	Geocoder *dataextraction.GeocodingService
	ESClient *storage.ElasticsearchClient
}

func NewProcessor(cfg *config.Config, geocoder *dataextraction.GeocodingService, esClient *storage.ElasticsearchClient) *Processor {
	return &Processor{
		Cfg:      cfg,
		Geocoder: geocoder,
		ESClient: esClient,
	}
}

func (p *Processor) Process(ctx context.Context, event models.RawOsintEvent) {
	log.Printf("Processing event: %s\n", string(event.Text))
	result, err := dataextraction.CreateAgentSummary(event, ctx, p.Cfg.OpenAIKey, p.Cfg.OpenAIModel)

	if err != nil {
		log.Printf("Error extracting info: %v\n", err)
		return
	}

	log.Printf("AI Summary: %+v\n", result)

	locationMap, geocodeErr := p.Geocoder.GetBatchCoordinates(result.EnLocations)
	if geocodeErr != nil {
		log.Printf("Error fetching coordinates: %v\n", geocodeErr)
		return
	}

	for loc, geo := range locationMap {
		err := p.ESClient.IndexGeocode(ctx, p.Cfg.ElasticsearchGeocodeIndex, loc, geo)
		if err != nil {
			log.Printf("Error indexing geocode for %s: %v\n", loc, err)
		}
	}

	processedEvent := storage.ProcessedEvent{
		RawMessage: event.Text,
		Summary:    result.HeSummary,
		Locations:  locationMap,
		Timestamp:  event.Date,
	}

	err = p.ESClient.IndexEvent(ctx, p.Cfg.ElasticsearchIndex, processedEvent)
	if err != nil {
		log.Printf("Error indexing event to elasticsearch: %v\n", err)
	} else {
		log.Println("Successfully indexed event to elasticsearch")
	}
}

func (p *Processor) StartWorker(ctx context.Context, taskQueue <-chan models.RawOsintEvent) {
	for event := range taskQueue {
		p.Process(ctx, event)
	}
}
