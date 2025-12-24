package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manos/favourites/favourites/application"
	"github.com/manos/favourites/favourites/domain"
	"github.com/manos/favourites/http/auth"
)

type favRepoTestStub struct {
	pageResp []domain.FavouriteEntity
	total    int
	updateCalled bool
}

func (r *favRepoTestStub) FindByUser(_ context.Context, _ string) ([]domain.FavouriteEntity, error) {
	return r.pageResp, nil
}
func (r *favRepoTestStub) FindByUserPage(_ context.Context, _ string, _ int, _ int) ([]domain.FavouriteEntity, int, error) {
	return r.pageResp, r.total, nil
}
func (r *favRepoTestStub) Save(_ context.Context, _ domain.FavouriteEntity) error { return nil }
func (r *favRepoTestStub) FindByID(_ context.Context, _ string) (domain.FavouriteEntity, error) {
	return domain.FavouriteEntity{ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight}, nil
}
func (r *favRepoTestStub) Delete(_ context.Context, _ string) error { return nil }
func (r *favRepoTestStub) UpdateDescription(_ context.Context, _ string, _ string) error {
	r.updateCalled = true
	return nil
}

type assetClientTestStub struct {
	assets map[string]application.AssetDTO
}

func (c *assetClientTestStub) GetInsight(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetAudience(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetChart(_ context.Context, _ string, _ string) (application.AssetDTO, error) {
	return application.AssetDTO{}, nil
}
func (c *assetClientTestStub) GetInsights(_ context.Context, _ string, ids []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}
func (c *assetClientTestStub) GetAudiences(_ context.Context, _ string, ids []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}
func (c *assetClientTestStub) GetCharts(_ context.Context, _ string, ids []string) (map[string]application.AssetDTO, error) {
	return c.assets, nil
}

func TestFavouritesHandler_ListWithPagination(t *testing.T) {
	repo := &favRepoTestStub{
		pageResp: []domain.FavouriteEntity{
			{ID: "fav-1", UserID: "user-1", AssetID: "ins-1", Type: domain.FavouriteInsight},
		},
		total: 1,
	}
	assets := &assetClientTestStub{
		assets: map[string]application.AssetDTO{"ins-1": {ID: "ins-1", Name: "n", Type: "INSIGHT"}},
	}
	svc := application.NewFavouriteService(repo, assets)
	h := NewFavouritesHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/favourites?page=2&limit=1", nil)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	h.HandleFavourites(rec, req.WithContext(ctx))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp application.FavouritePageDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Page != 2 || resp.Limit != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestFavouritesHandler_UpdateDescription(t *testing.T) {
	repo := &favRepoTestStub{
		pageResp: []domain.FavouriteEntity{},
		total:    0,
	}
	assets := &assetClientTestStub{}
	svc := application.NewFavouriteService(repo, assets)
	h := NewFavouritesHandler(svc)

	body := bytes.NewBufferString(`{"description":"new desc"}`)
	req := httptest.NewRequest(http.MethodPatch, "/favourites/fav-1", body)
	ctx := context.WithValue(req.Context(), auth.ContextKeySubject, "user-1")
	rec := httptest.NewRecorder()

	h.HandleFavourites(rec, req.WithContext(ctx))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if !repo.updateCalled {
		t.Fatalf("expected update description to be called")
	}
}
