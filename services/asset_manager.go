package services

import (
	"fmt"
	"github.com/google/uuid"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
)

type AssetManager struct {
	assetRepository repos.AssetRepository
	fileUploader    infra.FileUploader
}

func NewAssetManager(assetRepo repos.AssetRepository, fileUploader infra.FileUploader) IAssetManager {
	return &AssetManager{
		assetRepository: assetRepo,
		fileUploader:    fileUploader,
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

func (manager *AssetManager) UploadAsset(file *models.File, name string) (*models.Asset, error) {
	url, err := manager.fileUploader.UploadFile(file)
	if err != nil {
		return nil, err
	}
	return manager.NewAsset(name, file.MimeType, url)
}
