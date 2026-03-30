package nominatimgeocoder

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	de "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	"golang.org/x/time/rate"
)

func NominatimSearch(locationName string, limiter *rate.Limiter) (de.Geocode, *de.GeocodeError) {
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&addressdetails=1", url.QueryEscape(locationName))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeInvalidRequest, "failed to create request", err)
	}

	req.Header.Set("User-Agent", "OsintProcessingService/1.0 (ishaisela@gmail.com)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeNetworkError, "failed to execute request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeNetworkError, fmt.Sprintf("request failed with status code %d", resp.StatusCode), nil)
	}

	var apiResults nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeParsingError, "failed to decode response", err)
	}

	if len(apiResults) == 0 {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeNotFound, "no results found", nil)
	}

	placeRank := PlaceRank(apiResults[0].PlaceRank)

	log.Printf("Nominatim response for %s. PlaceRank %d, Code %s", locationName, placeRank, apiResults[0].Address.CountryCode)

	// Filter wide response like egypt
	if placeRank.IsWideScope() {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeFiltered, "result is too broad", nil)
	}

	if apiResults[0].Address.CountryCode != "il" {
		return de.Geocode{}, de.NewGeocodeError(de.ErrCodeFiltered, fmt.Sprintf("result is outside target country: %s", apiResults[0].Address.CountryCode), nil)
	}

	return de.Geocode{Lat: apiResults[0].Lat, Lon: apiResults[0].Lon}, nil
}
