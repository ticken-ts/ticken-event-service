package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-event-service/config"
	"ticken-event-service/infra"
)

type Factory struct {
	dbClient *mongo.Client
	dbName   string
}

func NewFactory(db infra.Db, dbConfig *config.DatabaseConfig) *Factory {
	return &Factory{
		dbClient: db.GetClient().(*mongo.Client),
		dbName:   dbConfig.Name,
	}
}

func (factory *Factory) BuildEventRepository() any {
	return NewEventRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildOrganizerRepository() any {
	return NewOrganizerRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildOrganizationRepository() any {
	return NewOrganizationRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildAssetRepository() any {
	return NewAssetRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildValidatorRepository() any {
	return NewValidatorRepository(factory.dbClient, factory.dbName)
}
