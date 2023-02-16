package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type AssetManager struct {
	assetRepository repos.AssetRepository
}

func NewAssetManager(assetRepo repos.AssetRepository) IAssetManager {
	return &AssetManager{
		assetRepository: assetRepo,
	}
}

func (manager *AssetManager) GetAsset(assetID uuid.UUID) (*models.Asset, error) {
	asset := manager.assetRepository.FindByID(assetID)
	if asset == nil {
		return nil, fmt.Errorf("asset with id %s not found", assetID.String())
	}
	return asset, nil
}

func (manager *AssetManager) NewAsset(name string, mimeType string, url string) (*models.Asset, error) {
	newAsset := models.NewAsset(uuid.New(), name, mimeType, url)
	err := manager.assetRepository.AddAsset(newAsset)
	if err != nil {
		return nil, err
	}
	return newAsset, nil
}
