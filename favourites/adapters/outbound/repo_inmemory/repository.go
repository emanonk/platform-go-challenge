package repo_inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/manos/favourites/favourites/application"
	"github.com/manos/favourites/favourites/domain"
)

type FavouriteRepository struct {
	mu   sync.RWMutex
	data map[string]domain.FavouriteEntity
}

func NewInMemoryFavouriteRepository() *FavouriteRepository {
	return &FavouriteRepository{
		data: map[string]domain.FavouriteEntity{
			"fav-1": {ID: "fav-1", UserID: "user-1", AssetID: "ins-001", Type: domain.FavouriteInsight},
			"fav-2": {ID: "fav-2", UserID: "user-1", AssetID: "aud-001", Type: domain.FavouriteAudience},
			"fav-3": {ID: "fav-3", UserID: "user-1", AssetID: "chr-001", Type: domain.FavouriteChart},
			"fav-4": {ID: "fav-4", UserID: "user-2", AssetID: "ins-002", Type: domain.FavouriteInsight},
		},
	}
}

func (r *FavouriteRepository) FindByUser(_ context.Context, userID string) ([]domain.FavouriteEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.FavouriteEntity, 0)
	for _, f := range r.data {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	log.Printf("repo favourites: find by user=%s count=%d", userID, len(out))
	return out, nil
}

func (r *FavouriteRepository) Save(_ context.Context, fav domain.FavouriteEntity) error {
	r.mu.Lock()
	r.data[fav.ID] = fav
	r.mu.Unlock()
	log.Printf("repo favourites: saved favourite_id=%s user=%s type=%s", fav.ID, fav.UserID, fav.Type)
	return nil
}

func (r *FavouriteRepository) FindByID(_ context.Context, favouriteID string) (domain.FavouriteEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.data[favouriteID]
	if ok {
		log.Printf("repo favourites: find by id=%s hit user=%s", favouriteID, f.UserID)
		return f, nil
	}
	log.Printf("repo favourites: find by id=%s not found", favouriteID)
	return domain.FavouriteEntity{}, application.ErrFavouriteNotFound
}

func (r *FavouriteRepository) Delete(_ context.Context, favouriteID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f, ok := r.data[favouriteID]; ok {
		delete(r.data, favouriteID)
		r.mu.Unlock()
		log.Printf("repo favourites: deleted id=%s user=%s", favouriteID, f.UserID)
		return nil
	}
	log.Printf("repo favourites: delete id=%s not found", favouriteID)
	return application.ErrFavouriteNotFound
}

func (r *FavouriteRepository) UpdateDescription(_ context.Context, favouriteID string, description string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f, ok := r.data[favouriteID]; ok {
		f.Description = description
		r.data[favouriteID] = f
		log.Printf("repo favourites: updated description id=%s user=%s", favouriteID, f.UserID)
		return nil
	}
	log.Printf("repo favourites: update description id=%s not found", favouriteID)
	return application.ErrFavouriteNotFound
}
