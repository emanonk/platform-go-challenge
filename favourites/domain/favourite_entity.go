package domain

type FavouriteType string

const (
	FavouriteInsight  FavouriteType = "INSIGHT"
	FavouriteAudience FavouriteType = "AUDIENCE"
	FavouriteChart    FavouriteType = "CHART"
)

type FavouriteEntity struct {
	ID      string
	UserID  string
	AssetID string
	Type    FavouriteType
}
