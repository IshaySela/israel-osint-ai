package models

import "strconv"

type GeocodeCache struct {
	LocationText string  `json:"location_text"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Timestamp    string  `json:"timestamp"`
}

func (geoCache *GeocodeCache) ToGeocode() Geocode {
	return Geocode{
		Lat: strconv.FormatFloat(geoCache.Lat, 'f', -1, 64),
		Lon: strconv.FormatFloat(geoCache.Lon, 'f', -1, 64),
	}
}
