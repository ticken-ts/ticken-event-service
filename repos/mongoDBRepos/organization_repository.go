package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
)

const TicketCollectionName = "organizations"

type OrganizationMongoDBRepository struct {
	baseRepository
}

func NewOrganizationRepository(db *mongo.Client, database string) *OrganizationMongoDBRepository {
	return &OrganizationMongoDBRepository{
		baseRepository{
			dbClient:       db,
			dbName:         database,
			collectionName: TicketCollectionName,
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

func (r *OrganizationMongoDBRepository) FindUserOrganization(userId string) *models.Organization {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizations := r.getCollection()
	org := organizations.FindOne(storeContext, bson.M{"users": userId})

	var foundOrg = new(models.Organization)
	err := org.Decode(foundOrg)
	if err != nil {
		return nil
	}
	return foundOrg
}

func (r *OrganizationMongoDBRepository) AddOrganization(org *models.Organization) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizations := r.getCollection()
	_, err := organizations.InsertOne(storeContext, org)
	if err != nil {
		return err
	}
	return nil
}
