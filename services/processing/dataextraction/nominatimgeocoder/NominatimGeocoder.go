package nominatimgeocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
)

func NominatimSearch(locationName string) (dataextraction.Geocode, *dataextraction.GeocodeError) {
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&addressdetails=1", url.QueryEscape(locationName))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeInvalidRequest, "failed to create request", err)
	}

	req.Header.Set("User-Agent", "OsintProcessingService/1.0 (ishaisela@gmail.com)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeNetworkError, "failed to execute request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeNetworkError, fmt.Sprintf("request failed with status code %d", resp.StatusCode), nil)
	}

	var apiResults geocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeParsingError, "failed to decode response", err)
	}

	if len(apiResults) == 0 {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeNotFound, "no results found", nil)
	}
	placeRank := PlaceRank(apiResults[0].PlaceRank)

	// Filter wide response like egypt
	if placeRank.IsWideScope() {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeFiltered, "result is too broad", nil)
	}

	if apiResults[0].Address.CountryCode != "il" {
		return dataextraction.Geocode{}, dataextraction.NewGeocodeError(dataextraction.ErrCodeFiltered, fmt.Sprintf("result is outside target country: %s", apiResults[0].Address.CountryCode), nil)
	}

	return dataextraction.Geocode{Lat: apiResults[0].Lat, Lon: apiResults[0].Lon}, nil
}
