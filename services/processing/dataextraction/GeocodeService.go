package dataextraction

import (
	"log"
	"sync"
)

type GeocoderFunction func(string) (Geocode, *GeocodeError)

type GeocodingService struct {
	mu       sync.RWMutex
	cache    map[string]Geocode
	geocoder GeocoderFunction
}

func NewGeocodingService(geocoder GeocoderFunction) *GeocodingService {
	return &GeocodingService{
		cache:    make(map[string]Geocode),
		geocoder: geocoder,
	}
}

func (s *GeocodingService) GetCoordinate(location string) (Geocode, *GeocodeError) {
	if location == "" {
		return Geocode{}, NewGeocodeError(ErrCodeInvalidRequest, "location string cannot be empty", nil)
	}

	s.mu.RLock()
	cached, exists := s.cache[location]
	s.mu.RUnlock()

	if exists {
		return cached, nil
	}

	geocode, err := s.geocoder(location)
	if err != nil {
		return Geocode{}, err
	}

	s.mu.Lock()
	s.cache[location] = geocode
	s.mu.Unlock()

	return geocode, nil
}

func (s *GeocodingService) GetBatchCoordinates(locations []string) (map[string]Geocode, *GeocodeError) {
	if len(locations) == 0 {
		return nil, NewGeocodeError(ErrCodeInvalidRequest, "locations list cannot be empty", nil)
	}

	results := make(map[string]Geocode)

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
