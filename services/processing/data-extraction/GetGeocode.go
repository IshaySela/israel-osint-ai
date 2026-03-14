package dataextraction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Geocode struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type GeocodeResponse []Geocode

type GeocodingService struct {
	mu    sync.RWMutex
	cache map[string]Geocode
}

func NewGeocodingService() *GeocodingService {
	return &GeocodingService{
		cache: make(map[string]Geocode),
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
		geocode, err := s.fetchFromAPI(location)
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

func (s *GeocodingService) fetchFromAPI(locationName string) (Geocode, error) {
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(locationName))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Geocode{}, err
	}

	req.Header.Set("User-Agent", "OsintProcessingService/1.0 (ishaisela@gmail.com)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Geocode{}, err
	}
	defer resp.Body.Close()

	var apiResults GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return Geocode{}, err
	}

	if len(apiResults) == 0 {
		return Geocode{}, fmt.Errorf("no results")
	}

	return apiResults[0], nil
}
