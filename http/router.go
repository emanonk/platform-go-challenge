package httpapi

import (
	"net/http"

	"github.com/manos/favourites/http/auth"
)

type Router struct {
	mux  *http.ServeMux
	root http.Handler
}

func NewRouter(
	jwtCfg auth.JWTConfig,
	favouritesHandler *FavouritesHandler,
	assetsHandler *AssetsHandler,
) *Router {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// protected /users/*
	protectedFavourites := Chain(
		http.HandlerFunc(favouritesHandler.HandleFavourites),
		auth.Middleware(jwtCfg),
		// auth.RequireUserMatch(extractUserIDFromJWT),
	)

	// protected /users/*
	protectedAssets := Chain(
		http.HandlerFunc(assetsHandler.GetAssetByID),
		auth.Middleware(jwtCfg),
		// auth.RequireUserMatch(extractUserIDFromJWT),
	)

	// mux.Handle("/users/", protectedFavourites)

	// favourites
	mux.Handle("/favourites", protectedFavourites)

	// assets
	mux.Handle("/assets", protectedAssets)

	return &Router{
		mux:  mux,
		root: mux,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func extractUserIDFromJWT(r *http.Request) (string, bool) {
	return auth.SubjectFromContext(r.Context())
}
