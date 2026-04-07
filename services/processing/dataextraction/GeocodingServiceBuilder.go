package dataextraction

import (
	"context"
	"fmt"

	"github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/geocodeerrors"
	nominatim "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/nominatimgeocoder"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
	"golang.org/x/time/rate"
)

type GeocodingServiceBuilder struct {
	ctx      context.Context
	geocoder GeocoderFunction
	cache    GeocodeCache
}

func NewGeocodingServiceBuilder() *GeocodingServiceBuilder {
	return &GeocodingServiceBuilder{}
}

func (b *GeocodingServiceBuilder) WithContext(ctx context.Context) *GeocodingServiceBuilder {
	b.ctx = ctx
	return b
}

func (b *GeocodingServiceBuilder) WithGeocodeFunction(fn GeocoderFunction) *GeocodingServiceBuilder {
	b.geocoder = fn
	return b
}

func (b *GeocodingServiceBuilder) WithNominatim(limiter *rate.Limiter) *GeocodingServiceBuilder {
	b.geocoder = func(location string) (models.Geocode, *geocodeerrors.GeocodeError) {
		return nominatim.NominatimSearch(location, limiter)
	}
	return b
}

func (b *GeocodingServiceBuilder) WithCache(cache GeocodeCache) *GeocodingServiceBuilder {
	b.cache = cache
	return b
}

func (b *GeocodingServiceBuilder) WithElasticsearchCache(es *storage.ElasticsearchClient) *GeocodingServiceBuilder {
	b.cache = es
	return b
}

func (b *GeocodingServiceBuilder) Build() (*GeocodingService, error) {
	if b.geocoder == nil {
		return nil, fmt.Errorf("geocoder function is required")
	}
	return &GeocodingService{
		ctx:          b.ctx,
		cache:        make(map[string]models.Geocode),
		geocoder:     b.geocoder,
		geocodeCache: b.cache,
	}, nil
}