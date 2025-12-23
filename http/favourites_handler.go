package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/manos/favourites/http/auth"

	"github.com/manos/favourites/favourites/application"
)

type FavouritesHandler struct {
	svc *application.FavouriteService
}

func NewFavouritesHandler(svc *application.FavouriteService) *FavouritesHandler {
	return &FavouritesHandler{svc: svc}
}

func (h *FavouritesHandler) HandleFavourites(w http.ResponseWriter, r *http.Request) {

	userId, _ := auth.SubjectFromContext(r.Context())
	// GET /favourites/
	// POST /favourites/
	// DELETE /favourites/{favouriteId}
	// PATCH /favourites/{favouriteId}

	fmt.Println("r.URL.Path:", r.URL.Path)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if parts[0] != "favourites" {
		http.NotFound(w, r)
		return
	}

	// fmt.Println("UserID in Handler:", userId)
	// fmt.Println("URL Parts:", parts)

	// userID := parts[1]
	fmt.Println("UserID from JWT:", userId, "Method:", r.Method)
	switch r.Method {
	case http.MethodGet:
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		h.listFavorites(w, r, userId)

	case http.MethodPost:
		// handle adding a favourite
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		h.AddFavourite(w, r, userId)
	case http.MethodDelete:
		// handle deleting a favourite
	case http.MethodPatch:
		// handle updating a favourite
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// result, err := h.svc.GetFavouritesForUser(r.Context(), userId)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(result)
}

func (h *FavouritesHandler) listFavorites(w http.ResponseWriter, r *http.Request, userId string) {
	favs, err := h.svc.GetFavouritesForUser(r.Context(), userId)

	fmt.Println("Favourites Result:", favs, err)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, favs)
}

func (h *FavouritesHandler) AddFavourite(w http.ResponseWriter, r *http.Request, userId string) {
	var newFavourite application.AddFavouriteRequest

	if err := json.NewDecoder(r.Body).Decode(&newFavourite); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Add Favourite Request:", newFavourite)
	fav, err := h.svc.AddFavourite(r.Context(), userId, newFavourite.Type, newFavourite.AssetID, newFavourite.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, fav)
}
