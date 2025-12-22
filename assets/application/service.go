package application

import (
	"context"
	"errors"

	"github.com/manos/favourites/assets/domain"
)

type AssetService struct {
	repo AssetRepository
}

func NewAssetService(repo AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

var ErrAssetNotFound = errors.New("asset not found")

// --- Facade ---

type AssetsService interface {
	GetInsight(ctx context.Context, assetID string) (domain.InsightAsset, error)
	GetAudience(ctx context.Context, assetID string) (domain.AudienceAsset, error)
	GetChart(ctx context.Context, assetID string) (domain.ChartAsset, error)
}

func (f *AssetService) GetInsight(ctx context.Context, assetID string) (domain.InsightAsset, error) {
	a, ok := f.repo.GetInsight(ctx, assetID)
	if !ok {
		return domain.InsightAsset{}, ErrAssetNotFound
	}
	return a, nil
}

func (f *AssetService) GetAudience(ctx context.Context, assetID string) (domain.AudienceAsset, error) {
	a, ok := f.repo.GetAudience(ctx, assetID)
	if !ok {
		return domain.AudienceAsset{}, ErrAssetNotFound
	}
	return a, nil
}

func (f *AssetService) GetChart(ctx context.Context, assetID string) (domain.ChartAsset, error) {
	a, ok := f.repo.GetChart(ctx, assetID)
	if !ok {
		return domain.ChartAsset{}, ErrAssetNotFound
	}
	return a, nil
}
