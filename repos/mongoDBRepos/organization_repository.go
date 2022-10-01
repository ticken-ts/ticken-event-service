package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
