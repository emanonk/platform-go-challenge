package application

import (
	"context"

	"github.com/manos/favourites/favourites/domain"
)

type FavouriteRepository interface {
	FindByUser(ctx context.Context, userID string) ([]domain.FavouriteEntity, error)
}

// Outbound port: Favourites depends on assets through this interface.
// Return favourites-owned DTOs (anti-corruption layer).
type AssetClient interface {
	GetInsight(ctx context.Context, assetID string) (AssetDTO, error)
	GetAudience(ctx context.Context, assetID string) (AssetDTO, error)
	GetChart(ctx context.Context, assetID string) (AssetDTO, error)
}
