package dataextraction

import (
	"context"
	"log"

	"processing/dataextraction/geocodeerrors"
	models "processing/models"
)

type GeocoderFunction func(string) (models.Geocode, *geocodeerrors.GeocodeError)

// GeocodeCache is a persistent cache for geocode results.
// Implementations must be safe for concurrent use by multiple goroutines.
type GeocodeCache interface {
	IndexGeocode(ctx context.Context, locationText string, geocode models.Geocode) (error, string)
	GetGeocode(ctx context.Context, location string) (models.GeocodeCache, error)
}

type GeocodingService struct {
	ctx          context.Context
	geocoder     GeocoderFunction
	geocodeCache GeocodeCache
}

func (s *GeocodingService) GetCoordinate(location string) (models.Geocode, *geocodeerrors.GeocodeError) {
	if location == "" {
		return models.Geocode{}, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeInvalidRequest, "location string cannot be empty", nil)
	}

	cached, err := s.geocodeCache.GetGeocode(s.ctx, location)

	if err == nil {
		return cached.ToGeocode(), nil
	}

	geocode, gecErr := s.geocoder(location)
	if gecErr != nil {
		return models.Geocode{}, gecErr
	}

	err, _ = s.geocodeCache.IndexGeocode(s.ctx, location, geocode)

	if err != nil {
		log.Printf("Error while caching result of geocode (%s): %v", location, err)
	}

	return geocode, nil
}

func (s *GeocodingService) GetBatchCoordinates(locations []string) (map[string]models.Geocode, *geocodeerrors.GeocodeError) {
	if len(locations) == 0 {
		return nil, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeInvalidRequest, "locations list cannot be empty", nil)
	}

	results := make(map[string]models.Geocode)

	for _, location := range locations {

		geocode, gecErr := s.GetCoordinate(location)
		if gecErr != nil {
			log.Printf("Warning: failed to fetch %s: %v\n", location, gecErr)
			continue
		}

		results[location] = geocode
	}

	if len(results) == 0 {
		return nil, geocodeerrors.NewGeocodeError(geocodeerrors.ErrCodeNotFound, "no locations found", nil)
	}

	return results, nil
}
