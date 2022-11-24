package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
)

const OrganizationCollectionName = "organizations"

type OrganizationMongoDBRepository struct {
	baseRepository
}

func NewOrganizationRepository(dbClient *mongo.Client, dbName string) *OrganizationMongoDBRepository {
	return &OrganizationMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: OrganizationCollectionName,
		},
	}
}

func (r *OrganizationMongoDBRepository) getCollection() *mongo.Collection {
	ctx, cancel := r.generateOpSubcontext()
	defer cancel()

	coll := r.baseRepository.getCollection()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "organization_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("error creating event repo: " + err.Error())
	}
	return coll
}

func (r *OrganizationMongoDBRepository) AddOrganization(organization *models.Organization) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	_, err := organizers.InsertOne(storeContext, organization)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrganizationMongoDBRepository) FindOrganization(organizerID string) *models.Organization {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"organization_id": organizerID})

	var foundOrganization models.Organization
	err := result.Decode(&foundOrganization)

	if err != nil {
		return nil
	}

	return &foundOrganization
}
