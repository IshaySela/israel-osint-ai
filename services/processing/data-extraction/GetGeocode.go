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

type nominatimGeocodeObject struct {
	PlaceID     int64    `json:"place_id"`
	OsmType     string   `json:"osm_type"`
	OsmID       int64    `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"` // e.g., "administrative" or "embassy"
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	AddressType string   `json:"addresstype"` // e.g., "country", "city"
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
}

type geocodeResponse []nominatimGeocodeObject

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
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&countrycodes=il", url.QueryEscape(locationName))

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

	var apiResults geocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return Geocode{}, err
	}

	if len(apiResults) == 0 {
		return Geocode{}, fmt.Errorf("no results")
	}
	placeRank := PlaceRank(apiResults[0].PlaceRank)

	// Filter wide response like egypt
	if placeRank.IsWideScope() {
		return Geocode{}, fmt.Errorf("no results")
	}

	return Geocode{Lat: apiResults[0].Lat, Lon: apiResults[0].Lon}, nil
}
