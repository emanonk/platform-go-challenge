package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/manos/favourites/favourites/application"
)

type FavouritesHandler struct {
	svc *application.FavouriteService
}

func NewFavouritesHandler(svc *application.FavouriteService) *FavouritesHandler {
	return &FavouritesHandler{svc: svc}
}

func (h *FavouritesHandler) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
	// GET /users/{userId}/favourites
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "users" || parts[2] != "favourites" {
		http.NotFound(w, r)
		return
	}

	userID := parts[1]

	result, err := h.svc.GetFavouritesForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
