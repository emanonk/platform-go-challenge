package application

import (
	"context"

	"github.com/manos/favourites/assets/domain"
)

// AssetRepository defines the minimal interface the application needs from a repository.
type AssetRepository interface {
	GetInsight(ctx context.Context, id string) (domain.InsightAsset, bool)
	GetAudience(ctx context.Context, id string) (domain.AudienceAsset, bool)
	GetChart(ctx context.Context, id string) (domain.ChartAsset, bool)
}
