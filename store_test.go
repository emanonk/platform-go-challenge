package main

import "testing"

func TestStoreAddAndList(t *testing.T) {
	store := NewStore()
	userID := "user-1"

	asset := Asset{
		Type:        AssetTypeInsight,
		Description: "My insight",
		Insight:     &InsightData{Text: "40% of millennials..."},
	}

	fav := store.AddFavorite(userID, asset)
	if fav.ID == "" || fav.Asset.ID == "" {
		t.Fatalf("expected IDs to be set")
	}

	list := store.ListFavorites(userID)
	if len(list) != 1 {
		t.Fatalf("expected 1 favourite, got %d", len(list))
	}
}

func TestStoreUpdateDescription(t *testing.T) {
	store := NewStore()
	userID := "user-1"

	asset := Asset{
		Type:        AssetTypeInsight,
		Description: "Old",
		Insight:     &InsightData{Text: "Something"},
	}
	fav := store.AddFavorite(userID, asset)

	updated, err := store.UpdateFavoriteDescription(userID, fav.ID, "New")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Asset.Description != "New" {
		t.Fatalf("expected description to be updated")
	}
}
