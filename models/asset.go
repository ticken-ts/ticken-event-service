package models

import "github.com/google/uuid"

type Asset struct {
	AssetID  uuid.UUID `bson:"asset_id"`
	Name     string    `bson:"name"`
	MimeType string    `bson:"mimeType"`
	URL      string    `bson:"url"`
}

func (asset *Asset) IsStoredLocally() bool {
	return len(asset.URL) >= 4 && asset.URL[0:4] != "http"
}
