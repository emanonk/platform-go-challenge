package application

import (
	"context"
	"fmt"

	"github.com/manos/favourites/favourites/domain"
)

type FavouriteService struct {
	repo   FavouriteRepository
	assets AssetClient
}

func NewFavouriteService(repo FavouriteRepository, assets AssetClient) *FavouriteService {
	return &FavouriteService{repo: repo, assets: assets}
}

func (s *FavouriteService) GetFavouritesForUser(ctx context.Context, userID string) ([]FavouriteDTO, error) {
	favs, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]FavouriteDTO, 0, len(favs))

	for _, fav := range favs {
		var asset AssetDTO

		switch fav.Type {
		case domain.FavouriteInsight:
			asset, err = s.assets.GetInsight(ctx, fav.AssetID)
		case domain.FavouriteAudience:
			asset, err = s.assets.GetAudience(ctx, fav.AssetID)
		case domain.FavouriteChart:
			asset, err = s.assets.GetChart(ctx, fav.AssetID)
		default:
			return nil, fmt.Errorf("unknown favourite type: %s", fav.Type)
		}

		if err != nil {
			return nil, err
		}

		out = append(out, FavouriteDTO{
			ID:     fav.ID,
			UserID: fav.UserID,
			Type:   string(fav.Type),
			Asset:  asset,
		})
	}

	return out, nil
}
