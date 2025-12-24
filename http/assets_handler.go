package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/manos/favourites/http/auth"

	application "github.com/manos/favourites/assets/application"
)

type AssetsHandler struct {
	svc *application.AssetService
}

func NewAssetsHandler(svc *application.AssetService) *AssetsHandler {
	return &AssetsHandler{svc: svc}
}

func (h *AssetsHandler) GetAssetByID(w http.ResponseWriter, r *http.Request) {

	userId, _ := auth.SubjectFromContext(r.Context())

	// GET /assets/{type}/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "assets" {
		http.NotFound(w, r)
		return
	}

	assetType := parts[1]
	assetID := parts[2]
	log.Printf("assets: %s %s user=%s type=%s id=%s", r.Method, r.URL.Path, userId, assetType, assetID)

	var (
		result any
		err    error
	)

	switch assetType {
	case "insights":
		result, err = h.svc.GetInsight(r.Context(), userId, assetID)
	case "audiences":
		result, err = h.svc.GetAudience(r.Context(), userId, assetID)
	case "charts":
		result, err = h.svc.GetChart(r.Context(), userId, assetID)
	default:
		log.Printf("assets: user=%s type=%s id=%s unknown type", userId, assetType, assetID)
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Printf("assets: user=%s type=%s id=%s err=%v", userId, assetType, assetID, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("assets: user=%s type=%s id=%s served", userId, assetType, assetID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
