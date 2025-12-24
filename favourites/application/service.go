package application

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/manos/favourites/favourites/domain"
)

var (
	ErrFavouriteNotFound = errors.New("favourite not found")
	ErrFavouriteForbidden = errors.New("favourite does not belong to user")
)

type FavouriteService struct {
	repo   FavouriteRepository
	assets AssetClient
}

func NewFavouriteService(repo FavouriteRepository, assets AssetClient) *FavouriteService {
	return &FavouriteService{repo: repo, assets: assets}
}

func (s *FavouriteService) GetFavouritesForUser(ctx context.Context, userId string) ([]FavouriteDTO, error) {
	log.Printf("svc favourites: list user=%s", userId)
	favs, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		log.Printf("svc favourites: list user=%s err=%v", userId, err)
		return nil, err
	}

	out := make([]FavouriteDTO, 0, len(favs))

	//todo change ask of list of assets ids to get all at once
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

	log.Printf("svc favourites: add user=%s type=%s asset=%s", userId, favType, assetId)
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
		log.Printf("svc favourites: add user=%s type=%s asset=%s validation err=%v", userId, favType, assetId, err)
		return "", err
	}

	fav, err := domain.NewFavourite(userId, domain.FavouriteType(favType), assetId, description)
	if err != nil {
		log.Printf("svc favourites: add user=%s type=%s asset=%s build err=%v", userId, favType, assetId, err)
		return "", err
	}
	err = s.repo.Save(ctx, fav)
	if err != nil {
		log.Printf("svc favourites: add user=%s type=%s asset=%s save err=%v", userId, favType, assetId, err)
		return "", err
	}
	log.Printf("svc favourites: add user=%s favourite_id=%s created", userId, fav.ID)
	return fav.ID, nil
}

func (s *FavouriteService) DeleteFavourite(ctx context.Context, userId string, favouriteID string) error {
	log.Printf("svc favourites: delete user=%s favourite_id=%s", userId, favouriteID)
	fav, err := s.repo.FindByID(ctx, favouriteID)
	if err != nil {
		log.Printf("svc favourites: delete user=%s favourite_id=%s find err=%v", userId, favouriteID, err)
		return err
	}

	if fav.UserID != userId {
		log.Printf("svc favourites: delete user=%s favourite_id=%s forbidden owner=%s", userId, favouriteID, fav.UserID)
		return ErrFavouriteForbidden
	}

	return s.repo.Delete(ctx, favouriteID)
}
