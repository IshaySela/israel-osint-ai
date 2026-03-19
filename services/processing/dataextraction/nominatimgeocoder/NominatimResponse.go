package nominatimgeocoder

type nominatimAddress struct {
	Road        string `json:"road,omitempty"`
	Suburb      string `json:"suburb,omitempty"`
	City        string `json:"city,omitempty"`
	Town        string `json:"town,omitempty"`
	State       string `json:"state,omitempty"`
	Postcode    string `json:"postcode,omitempty"`
	Country     string `json:"country"`      // Full country name
	CountryCode string `json:"country_code"` // ISO 3166-1 alpha-2 code
}

type nominatimGeocodeObject struct {
	PlaceID     int64            `json:"place_id"`
	OsmType     string           `json:"osm_type"`
	OsmID       int64            `json:"osm_id"`
	Lat         string           `json:"lat"`
	Lon         string           `json:"lon"`
	Class       string           `json:"class"`
	Type        string           `json:"type"` // e.g., "administrative" or "embassy"
	PlaceRank   int              `json:"place_rank"`
	Importance  float64          `json:"importance"`
	AddressType string           `json:"addresstype"` // e.g., "country", "city"
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	BoundingBox []string         `json:"boundingbox"`
	Address     nominatimAddress `json:"address"`
}

type geocodeResponse []nominatimGeocodeObject
