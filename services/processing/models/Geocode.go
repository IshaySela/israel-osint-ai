package models

type Geocode struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Location struct {
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}
