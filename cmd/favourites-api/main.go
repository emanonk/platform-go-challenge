package main

import (
	"log"
	"net/http"
	"time"

	assetrepo "github.com/manos/favourites/assets/adapters/outbound/repo_inmemory"
	assetapplication "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/favourites/adapters/outbound/assets_client"

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

	router := httpapi.NewRouter(favHandler, assetsHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
