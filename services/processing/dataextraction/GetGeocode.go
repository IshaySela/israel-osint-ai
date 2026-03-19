package dataextraction

import (
	"log"
	"sync"
	"time"
)

type GeocoderFunction func(string) (Geocode, error)

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

func (s *GeocodingService) GetBatchCoordinates(locations []string) ([]Geocode, error) {
	results := make([]Geocode, 0, len(locations))
	ticker := time.NewTicker(1100 * time.Millisecond)
	defer ticker.Stop()

	for _, location := range locations {
		s.mu.RLock()
		cached, exists := s.cache[location]
		s.mu.RUnlock()

		if exists {
			results = append(results, cached)
			continue
		}

		<-ticker.C
		geocode, err := s.geocoder(location)
		if err != nil {
			log.Printf("Warning: failed to fetch %s: %v\n", location, err)
			continue
		}

		s.mu.Lock()
		s.cache[location] = geocode
		s.mu.Unlock()

		results = append(results, geocode)
	}

	return results, nil
}
