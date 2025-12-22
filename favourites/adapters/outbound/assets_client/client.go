package assets_client

import (
	"context"

	assetsapp "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/favourites/application"
)

type Client struct {
	assetService assetsapp.AssetsService
}

func NewAssetClient(assetService assetsapp.AssetsService) *Client {
	return &Client{assetService: assetService}
}

func (c *Client) GetInsight(ctx context.Context, assetID string) (application.AssetDTO, error) {
	a, err := c.assetService.GetInsight(ctx, assetID)
	if err != nil {
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		OwnerUserID: a.UserId,
		Type:        "INSIGHT",
		Text:        a.Text,
	}, nil
}

func (c *Client) GetAudience(ctx context.Context, assetID string) (application.AssetDTO, error) {
	a, err := c.assetService.GetAudience(ctx, assetID)
	if err != nil {
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:                a.Id,
		Name:              a.Name,
		Description:       a.Description,
		OwnerUserID:       a.UserId,
		Type:              "AUDIENCE",
		SampleSize:        a.SampleSize,
		TotalRespondents:  a.TotalRespondents,
		EstimatedReach:    a.EstimatedReach,
		PopulationPercent: a.PopulationPercent,
	}, nil
}

func (c *Client) GetChart(ctx context.Context, assetID string) (application.AssetDTO, error) {
	a, err := c.assetService.GetChart(ctx, assetID)
	if err != nil {
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		OwnerUserID: a.UserId,
		Type:        "CHART",
		XAxisTitle:  a.XAxisTitle,
		YAxisTitle:  a.YAxisTitle,
		Data:        a.Data,
	}, nil
}
