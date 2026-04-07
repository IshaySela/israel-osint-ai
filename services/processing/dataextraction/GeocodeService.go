package dataextraction

import (
	"context"
	"log"
	"sync"

	"github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/geocodeerrors"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

type GeocoderFunction func(string) (models.Geocode, *geocodeerrors.GeocodeError)

// GeocodeCache is a persistent cache for geocode results.
// Implementations must be safe for concurrent use by multiple goroutines.
type GeocodeCache interface {
	IndexGeocode(ctx context.Context, locationText string, geocode models.Geocode) (error, string)
	GetGeocode(ctx context.Context, index string, location string) (models.GeocodeCache, error)
}

type GeocodingService struct {
	mu           sync.RWMutex
	ctx          context.Context
	cache        map[string]models.Geocode
	geocoder     GeocoderFunction
	geocodeCache GeocodeCache
}

func (s *GeocodingService) GetCoordinate(location string) (models.Geocode, *geocodeerrors.GeocodeError) {
	if location == "" {
		return models.Geocode{}, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeInvalidRequest, "location string cannot be empty", nil)
	}

	s.mu.RLock()
	cached, exists := s.cache[location]
	s.mu.RUnlock()

	if exists {
		return cached, nil
	}

	geocode, err := s.geocoder(location)
	if err != nil {
		return models.Geocode{}, err
	}

	s.mu.Lock()
	s.cache[location] = geocode
	s.mu.Unlock()

	return geocode, nil
}

func (s *GeocodingService) GetBatchCoordinates(locations []string) (map[string]models.Geocode, *geocodeerrors.GeocodeError) {
	if len(locations) == 0 {
		return nil, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeInvalidRequest, "locations list cannot be empty", nil)
	}

	results := make(map[string]models.Geocode)

	for _, location := range locations {
		if location == "" {
			continue
		}

		s.mu.RLock()
		cached, exists := s.cache[location]
		s.mu.RUnlock()

		if exists {
			results[location] = cached
			continue
		}

		geocode, err := s.geocoder(location)
		if err != nil {
			log.Printf("Warning: failed to fetch %s: %v\n", location, err)
			continue
		}

		s.mu.Lock()
		s.cache[location] = geocode
		s.mu.Unlock()

		results[location] = geocode
	}

	if len(results) == 0 {
		return nil, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeNotFound, "no locations found", nil)
	}

	return results, nil
}
