package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	log.Printf("favourites: %s %s user=%s", r.Method, r.URL.Path, userId)
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
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		h.deleteFavourite(w, r, userId, parts[1])
	case http.MethodPatch:
		// handle updating a favourite
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		h.updateFavourite(w, r, userId, parts[1])
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

	if err != nil {
		log.Printf("favourites list: user=%s err=%v", userId, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	log.Printf("favourites list: user=%s count=%d", userId, len(favs))
	writeJSON(w, http.StatusOK, favs)
}

func (h *FavouritesHandler) AddFavourite(w http.ResponseWriter, r *http.Request, userId string) {
	var newFavourite application.AddFavouriteRequest

	if err := json.NewDecoder(r.Body).Decode(&newFavourite); err != nil {
		log.Printf("favourites add: user=%s decode error: %v", userId, err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("favourites add: user=%s type=%s asset=%s", userId, newFavourite.Type, newFavourite.AssetID)
	fav, err := h.svc.AddFavourite(r.Context(), userId, newFavourite.Type, newFavourite.AssetID, newFavourite.Description)
	if err != nil {
		log.Printf("favourites add: user=%s type=%s asset=%s err=%v", userId, newFavourite.Type, newFavourite.AssetID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("favourites add: user=%s favourite_id=%s created", userId, fav)
	writeJSON(w, http.StatusCreated, fav)
}

func (h *FavouritesHandler) deleteFavourite(w http.ResponseWriter, r *http.Request, userId, favouriteID string) {
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

type updateFavouriteRequest struct {
	Description string `json:"description"`
}

func (h *FavouritesHandler) updateFavourite(w http.ResponseWriter, r *http.Request, userId, favouriteID string) {
	var req updateFavouriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("favourites update: user=%s favourite_id=%s decode error: %v", userId, favouriteID, err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("favourites update: user=%s favourite_id=%s", userId, favouriteID)
	if err := h.svc.UpdateFavouriteDescription(r.Context(), userId, favouriteID, req.Description); err != nil {
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
