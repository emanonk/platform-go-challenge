package repo_inmemory

import (
	"context"

	"github.com/manos/favourites/favourites/application"
	"github.com/manos/favourites/favourites/domain"
)

// todo  change to map and add mutex for concurrency in real app
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

func (r *FavouriteRepository) Save(_ context.Context, fav domain.FavouriteEntity) error {
	r.data = append(r.data, fav)
	return nil
}

func (r *FavouriteRepository) FindByID(_ context.Context, favouriteID string) (domain.FavouriteEntity, error) {
	for _, f := range r.data {
		if f.ID == favouriteID {
			return f, nil
		}
	}
	return domain.FavouriteEntity{}, application.ErrFavouriteNotFound
}

func (r *FavouriteRepository) Delete(_ context.Context, favouriteID string) error {
	for i, f := range r.data {
		if f.ID == favouriteID {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return application.ErrFavouriteNotFound
}
