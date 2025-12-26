package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/manos/favourites/http/auth"

	"github.com/manos/favourites/favourites/application"
)

type FavouritesHandler struct {
	svc          *application.FavouriteService
	defaultPage  int
	defaultLimit int
	maxLimit     int
}

func NewFavouritesHandler(svc *application.FavouriteService, defaultPage int, defaultLimit int, maxLimit int) *FavouritesHandler {
	if defaultPage <= 0 {
		defaultPage = 1
	}
	if defaultLimit <= 0 {
		defaultLimit = 20
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}
	if defaultLimit > maxLimit {
		defaultLimit = maxLimit
	}
	return &FavouritesHandler{
		svc:          svc,
		defaultPage:  defaultPage,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

func (h *FavouritesHandler) GetFavourites(w http.ResponseWriter, r *http.Request, params GetFavouritesParams) {
	userId, _ := auth.SubjectFromContext(r.Context())
	page := h.defaultPage
	limit := h.defaultLimit

	if params.Page != nil {
		page = *params.Page
		if page < 1 {
			http.Error(w, "page must be >= 1", http.StatusBadRequest)
			return
		}
	}
	if params.Limit != nil {
		limit = *params.Limit
		if limit < 1 {
			http.Error(w, "limit must be >= 1", http.StatusBadRequest)
			return
		}
	}

	if limit > h.maxLimit {
		limit = h.maxLimit
	}

	favs, err := h.svc.GetFavouritesForUser(r.Context(), userId, page, limit)
	if err != nil {
		log.Printf("favourites list: user=%s err=%v", userId, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	log.Printf("favourites list: user=%s page=%d limit=%d count=%d total=%d", userId, page, limit, len(favs.Items), favs.Total)

	writeJSON(w, http.StatusOK, favouritePageToAPI(favs))
}

func (h *FavouritesHandler) PostFavourites(w http.ResponseWriter, r *http.Request) {
	userId, _ := auth.SubjectFromContext(r.Context())

	var req AddFavouriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("favourites add: user=%s decode error: %v", userId, err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	log.Printf("favourites add: user=%s type=%s asset=%s", userId, req.Type, req.AssetId)
	fav, err := h.svc.AddFavourite(r.Context(), userId, string(req.Type), req.AssetId, desc)
	if err != nil {
		log.Printf("favourites add: user=%s type=%s asset=%s err=%v", userId, req.Type, req.AssetId, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("favourites add: user=%s favourite_id=%s created", userId, fav)
	writeJSON(w, http.StatusCreated, fav)
}

func (h *FavouritesHandler) DeleteFavouritesId(w http.ResponseWriter, r *http.Request, favouriteID string) {
	userId, _ := auth.SubjectFromContext(r.Context())
	err := h.svc.DeleteFavourite(r.Context(), userId, favouriteID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrFavouriteNotFound):
			log.Printf("favourites delete: user=%s favourite_id=%s not found", userId, favouriteID)
			http.NotFound(w, r)
		case errors.Is(err, application.ErrFavouriteForbidden):
			log.Printf("favourites delete: user=%s favourite_id=%s forbidden", userId, favouriteID)
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			log.Printf("favourites delete: user=%s favourite_id=%s err=%v", userId, favouriteID, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("favourites delete: user=%s favourite_id=%s deleted", userId, favouriteID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouritesHandler) PatchFavouritesId(w http.ResponseWriter, r *http.Request, favouriteID string) {
	userId, _ := auth.SubjectFromContext(r.Context())

	var req UpdateFavouriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("favourites update: user=%s favourite_id=%s decode error: %v", userId, favouriteID, err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	log.Printf("favourites update: user=%s favourite_id=%s", userId, favouriteID)
	if err := h.svc.UpdateFavouriteDescription(r.Context(), userId, favouriteID, desc); err != nil {
		switch {
		case errors.Is(err, application.ErrFavouriteNotFound):
			log.Printf("favourites update: user=%s favourite_id=%s not found", userId, favouriteID)
			http.NotFound(w, r)
		case errors.Is(err, application.ErrFavouriteForbidden):
			log.Printf("favourites update: user=%s favourite_id=%s forbidden", userId, favouriteID)
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			log.Printf("favourites update: user=%s favourite_id=%s err=%v", userId, favouriteID, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func favouritePageToAPI(dto application.FavouritePageDTO) FavouritePage {
	items := make([]Favourite, 0, len(dto.Items))
	for _, f := range dto.Items {
		items = append(items, favouriteToAPI(f))
	}

	return FavouritePage{
		Items:      &items,
		Page:       intPtr(dto.Page),
		Limit:      intPtr(dto.Limit),
		Total:      intPtr(dto.Total),
		TotalPages: intPtr(dto.TotalPages),
	}
}

func favouriteToAPI(f application.FavouriteDTO) Favourite {
	typ := FavouriteType(f.Type)
	return Favourite{
		Id:          stringPtr(f.ID),
		UserId:      stringPtr(f.UserID),
		Type:        &typ,
		Description: stringPtr(f.Description),
		Asset:       assetToAPI(f.Asset),
	}
}

func assetToAPI(a application.AssetDTO) *Asset {
	typ := AssetType(a.Type)
	out := &Asset{
		Id:          stringPtr(a.ID),
		Name:        stringPtr(a.Name),
		Description: stringPtr(a.Description),
		OwnerUserId: stringPtr(a.OwnerUserID),
		Type:        &typ,
	}

	if a.Text != "" {
		out.Text = stringPtr(a.Text)
	}
	if a.XAxisTitle != "" {
		out.XAxisTitle = stringPtr(a.XAxisTitle)
	}
	if a.YAxisTitle != "" {
		out.YAxisTitle = stringPtr(a.YAxisTitle)
	}
	if len(a.Data) > 0 {
		data := make([]float32, 0, len(a.Data))
		for _, v := range a.Data {
			data = append(data, float32(v))
		}
		out.Data = &data
	}
	if a.SampleSize != 0 {
		out.SampleSize = intPtr(int(a.SampleSize))
	}
	if a.TotalRespondents != 0 {
		out.TotalRespondents = intPtr(int(a.TotalRespondents))
	}
	if a.EstimatedReach != 0 {
		out.EstimatedReach = intPtr(int(a.EstimatedReach))
	}
	if a.PopulationPercent != 0 {
		val := float32(a.PopulationPercent)
		out.PopulationPercent = &val
	}

	return out
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}
