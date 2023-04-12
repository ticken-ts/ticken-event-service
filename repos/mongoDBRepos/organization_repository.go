package mongoDBRepos

import (
	"ticken-event-service/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			primKeyField:   "organization_id",
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

func (r *OrganizationMongoDBRepository) FindOrganization(organizerID uuid.UUID) *models.Organization {
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

func (r *OrganizationMongoDBRepository) FindByName(name string) *models.Organization {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"name": name})

	var foundOrganization models.Organization
	err := result.Decode(&foundOrganization)

	if err != nil {
		return nil
	}

	return &foundOrganization
}

func (r *OrganizationMongoDBRepository) FindByMSPID(mspID string) *models.Organization {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"msp_id": mspID})

	var foundOrganization models.Organization
	err := result.Decode(&foundOrganization)

	if err != nil {
		return nil
	}

	return &foundOrganization
}

func (r *OrganizationMongoDBRepository) AnyWithName(name string) bool {
	return r.FindByName(name) != nil
}

func (r *OrganizationMongoDBRepository) FindByOrganizer(organizerID uuid.UUID) []*models.Organization {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizations := r.getCollection()

	result, err := organizations.Aggregate(findContext, mongo.Pipeline{
		// get orgs such that organizerID is the organizer_id of a user in the users list
		{{Key: "$match", Value: bson.M{"users": bson.M{"$elemMatch": bson.M{"organizer_id": organizerID}}}}},
	})

	if err != nil {
		return nil
	}

	var foundOrganizations []*models.Organization
	err = result.All(findContext, &foundOrganizations)

	if err != nil {
		return nil
	}

	return foundOrganizations
}
