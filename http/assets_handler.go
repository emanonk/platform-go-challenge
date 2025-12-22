package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	application "github.com/manos/favourites/assets/application"
)

type AssetsHandler struct {
	svc *application.AssetService
}

func NewAssetsHandler(svc *application.AssetService) *AssetsHandler {
	return &AssetsHandler{svc: svc}
}

func (h *AssetsHandler) GetAssetByID(w http.ResponseWriter, r *http.Request) {
	// GET /assets/{type}/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "assets" {
		http.NotFound(w, r)
		return
	}

	assetType := parts[1]
	assetID := parts[2]

	var (
		result any
		err    error
	)

	switch assetType {
	case "insights":
		result, err = h.svc.GetInsight(r.Context(), assetID)
	case "audiences":
		result, err = h.svc.GetAudience(r.Context(), assetID)
	case "charts":
		result, err = h.svc.GetChart(r.Context(), assetID)
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
