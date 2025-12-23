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

func (s *FavouriteService) GetFavouritesForUser(ctx context.Context, userId string) ([]FavouriteDTO, error) {
	favs, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	out := make([]FavouriteDTO, 0, len(favs))

	for _, fav := range favs {
		var asset AssetDTO

		switch fav.Type {
		case domain.FavouriteInsight:
			asset, err = s.assets.GetInsight(ctx, userId, fav.AssetID)
		case domain.FavouriteAudience:
			asset, err = s.assets.GetAudience(ctx, userId, fav.AssetID)
		case domain.FavouriteChart:
			asset, err = s.assets.GetChart(ctx, userId, fav.AssetID)
		default:
			return nil, fmt.Errorf("unknown favourite type: %s", fav.Type)
		}

		if err != nil {
			return nil, err
		}

		out = append(out, FavouriteDTO{
			ID:          fav.ID,
			UserID:      fav.UserID,
			Type:        string(fav.Type),
			Description: fav.Description,
			Asset:       asset,
		})
	}

	return out, nil
}

func (s *FavouriteService) AddFavourite(ctx context.Context, userId string, favType string, assetId string, description string) (string, error) {

	// var asset AssetDTO
	var err error

	favouriteType := domain.FavouriteType(favType)
	switch favouriteType {
	case domain.FavouriteInsight:
		_, err = s.assets.GetInsight(ctx, userId, assetId)
	case domain.FavouriteAudience:
		_, err = s.assets.GetAudience(ctx, userId, assetId)
	case domain.FavouriteChart:
		_, err = s.assets.GetChart(ctx, userId, assetId)
	default:
		return "", fmt.Errorf("unknown favourite type: %s", favouriteType)
	}

	if err != nil {
		return "", err
	}

	fav, err := domain.NewFavourite(userId, domain.FavouriteType(favType), assetId, description)
	if err != nil {
		return "", err
	}
	err = s.repo.Save(ctx, fav)
	if err != nil {
		return "", err
	}
	return fav.ID, nil
}
