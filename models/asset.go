package models

import "github.com/google/uuid"

type Asset struct {
	ID       uuid.UUID `bson:"asset_id"`
	Name     string    `bson:"name"`
	MimeType string    `bson:"mimeType"`
	URL      string    `bson:"url"`
}

func NewAsset(id uuid.UUID, name string, mimeType string, url string) *Asset {
	return &Asset{
		ID:       id,
		Name:     name,
		MimeType: mimeType,
		URL:      url,
	}
}
