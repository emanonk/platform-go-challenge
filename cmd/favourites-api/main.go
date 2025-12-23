package main

import (
	"log"
	"net/http"
	"time"

	assetrepo "github.com/manos/favourites/assets/adapters/outbound/repo_inmemory"
	assetapplication "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/favourites/adapters/outbound/assets_client"
	"github.com/manos/favourites/http/auth"

	favrepo "github.com/manos/favourites/favourites/adapters/outbound/repo_inmemory"
	favouritesapplication "github.com/manos/favourites/favourites/application"
	httpapi "github.com/manos/favourites/http"
)

func main() {
	// assets hexagon
	assetsRepo := assetrepo.NewInMemoryAssetRepository()
	assetsService := assetapplication.NewAssetService(assetsRepo)

	// favourites hexagon

	assetClient := assets_client.NewAssetClient(assetsService)
	favRepo := favrepo.NewInMemoryFavouriteRepository()
	favService := favouritesapplication.NewFavouriteService(favRepo, assetClient)

	// http handlers
	favHandler := httpapi.NewFavouritesHandler(favService)
	assetsHandler := httpapi.NewAssetsHandler(assetsService)

	pubKey, err := auth.LoadRSAPublicKey("public.pem")
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}

	jwtCfg := auth.JWTConfig{
		PublicKey: pubKey,
		Issuer:    "favourites-api",
		Audience:  "web",
	}

	router := httpapi.NewRouter(jwtCfg, favHandler, assetsHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
