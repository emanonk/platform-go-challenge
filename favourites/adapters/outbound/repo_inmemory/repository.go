package repo_inmemory

import (
	"context"

	"github.com/manos/favourites/favourites/domain"
)

type FavouriteRepository struct {
	data []domain.FavouriteEntity
}

func NewInMemoryFavouriteRepository() *FavouriteRepository {
	return &FavouriteRepository{
		data: []domain.FavouriteEntity{
			{ID: "fav-1", UserID: "user-1", AssetID: "ins-001", Type: domain.FavouriteInsight},
			{ID: "fav-2", UserID: "user-1", AssetID: "aud-001", Type: domain.FavouriteAudience},
			{ID: "fav-3", UserID: "user-1", AssetID: "chr-001", Type: domain.FavouriteChart},
			{ID: "fav-4", UserID: "user-2", AssetID: "ins-002", Type: domain.FavouriteInsight},
		},
	}
}

func (r *FavouriteRepository) FindByUser(_ context.Context, userID string) ([]domain.FavouriteEntity, error) {
	out := make([]domain.FavouriteEntity, 0)
	for _, f := range r.data {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	return out, nil
}
