package services

import (
	"github.com/google/uuid"
	"ticken-event-service/infra"
	"ticken-event-service/models"
	"ticken-event-service/repos"
	"ticken-event-service/tickenerr"
	"ticken-event-service/tickenerr/asseterr"
	"ticken-event-service/utils/file"
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

func (manager *AssetManager) GetAssetURL(assetID uuid.UUID) (string, error) {
	asset := manager.assetRepository.FindByID(assetID)
	if asset == nil {
		return "", tickenerr.New(asseterr.AssetNotFoundErrorCode)
	}

	return asset.URL, nil
}

func (manager *AssetManager) UploadAsset(file *file.File, name string) (*models.Asset, error) {
	url, err := manager.fileUploader.UploadFile(file)
	if err != nil {
		return nil, tickenerr.FromError(asseterr.FailedToUploadAsset, err)
	}

	newAsset := &models.Asset{
		ID:       uuid.New(),
		Name:     name,
		MimeType: file.MimeType,
		URL:      url,
	}

	if err := manager.assetRepository.AddAsset(newAsset); err != nil {
		return nil, tickenerr.FromError(asseterr.FailedToStoreAssetInDatabase, err)
	}

	return newAsset, nil
}
