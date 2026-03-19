package nominatimgeocoder

// PlaceRank represents the Nominatim spatial granularity level.
// https://nominatim.org/release-docs/latest/customize/Ranking/#search-rank
type PlaceRank int

const (
	RankOceanContinent        PlaceRank = 1
	RankCountry               PlaceRank = 4
	RankStateRegionProvince   PlaceRank = 5
	RankCounty                PlaceRank = 10
	RankCityMunicipality      PlaceRank = 13
	RankTownBorough           PlaceRank = 17
	RankVillageSuburb         PlaceRank = 19
	RankHamletNeighbourhood   PlaceRank = 20
	RankIsolatedDwellingBlock PlaceRank = 21
)

// IsWideScope identifies geometries that represent large administrative areas.
func (pr PlaceRank) IsWideScope() bool {
	return pr < RankCityMunicipality
}
