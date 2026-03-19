package nominatimgeocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
)

func NominatimSearch(locationName string) (dataextraction.Geocode, error) {
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&addressdetails=1", url.QueryEscape(locationName))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return dataextraction.Geocode{}, err
	}

	req.Header.Set("User-Agent", "OsintProcessingService/1.0 (ishaisela@gmail.com)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return dataextraction.Geocode{}, err
	}
	defer resp.Body.Close()

	var apiResults geocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return dataextraction.Geocode{}, err
	}

	if len(apiResults) == 0 {
		return dataextraction.Geocode{}, fmt.Errorf("no results")
	}
	placeRank := PlaceRank(apiResults[0].PlaceRank)

	// Filter wide response like egypt
	if placeRank.IsWideScope() || apiResults[0].Address.CountryCode != "il" {
		return dataextraction.Geocode{}, fmt.Errorf("no results")
	}

	return dataextraction.Geocode{Lat: apiResults[0].Lat, Lon: apiResults[0].Lon}, nil
}
