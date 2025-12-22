package httpapi

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func NewRouter(
	favouritesHandler *FavouritesHandler,
	assetsHandler *AssetsHandler,
) *Router {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// favourites
	mux.HandleFunc("/users/", favouritesHandler.GetUserFavourites)

	// assets
	mux.HandleFunc("/assets/", assetsHandler.GetAssetByID)

	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
