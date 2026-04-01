package dataextraction

import (
	"log"
	"sync"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

type GeocoderFunction func(string) (models.Geocode, *GeocodeError)

type GeocodingService struct {
	mu       sync.RWMutex
	cache    map[string]models.Geocode
	geocoder GeocoderFunction
}

func NewGeocodingService(geocoder GeocoderFunction) *GeocodingService {
	return &GeocodingService{
		cache:    make(map[string]models.Geocode),
		geocoder: geocoder,
	}
}

func (s *GeocodingService) GetCoordinate(location string) (models.Geocode, *GeocodeError) {
	if location == "" {
		return models.Geocode{}, NewGeocodeError(ErrCodeInvalidRequest, "location string cannot be empty", nil)
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

func (s *GeocodingService) GetBatchCoordinates(locations []string) (map[string]models.Geocode, *GeocodeError) {
	if len(locations) == 0 {
		return nil, NewGeocodeError(ErrCodeInvalidRequest, "locations list cannot be empty", nil)
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
		return nil, NewGeocodeError(ErrCodeNotFound, "no locations found", nil)
	}

	return results, nil
}
