package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
)

const AssetCollectionName = "assets"

type AssetMongoDBRepository struct {
	baseRepository
}

func NewAssetRepository(dbClient *mongo.Client, dbName string) *AssetMongoDBRepository {
	return &AssetMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: AssetCollectionName,
			primKeyField:   "asset_id",
		},
	}
}

func (r *AssetMongoDBRepository) getCollection() *mongo.Collection {
	ctx, cancel := r.generateOpSubcontext()
	defer cancel()

	coll := r.baseRepository.getCollection()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "asset_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("error creating event repo: " + err.Error())
	}
	return coll
}

func (r *AssetMongoDBRepository) AddAsset(asset *models.Asset) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	assets := r.getCollection()
	_, err := assets.InsertOne(storeContext, asset)
	if err != nil {
		return err
	}

	return nil
}

func (r *AssetMongoDBRepository) FindByID(assetID uuid.UUID) *models.Asset {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	assets := r.getCollection()
	result := assets.FindOne(findContext, bson.M{"asset_id": assetID})

	var foundAsset models.Asset
	err := result.Decode(&foundAsset)

	if err != nil {
		return nil
	}

	return &foundAsset
}
