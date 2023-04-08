package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
	"ticken-event-service/utils"
)

const OrganizerCollectionName = "organizers"

type OrganizerMongoDBRepository struct {
	baseRepository
}

func NewOrganizerRepository(dbClient *mongo.Client, dbName string) *OrganizerMongoDBRepository {
	return &OrganizerMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: OrganizerCollectionName,
			primKeyField:   utils.GetStructTag(models.Organizer{}.OrganizerID, "bson"),
		},
	}
}

func (r *OrganizerMongoDBRepository) getCollection() *mongo.Collection {
	ctx, cancel := r.generateOpSubcontext()
	defer cancel()

	coll := r.baseRepository.getCollection()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "organizer_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("error creating event repo: " + err.Error())
	}
	return coll
}

func (r *OrganizerMongoDBRepository) AddOrganizer(organizer *models.Organizer) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	_, err := organizers.InsertOne(storeContext, organizer)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrganizerMongoDBRepository) FindOrganizer(organizerID uuid.UUID) *models.Organizer {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"organizer_id": organizerID})

	var foundOrganizer models.Organizer
	err := result.Decode(&foundOrganizer)

	if err != nil {
		return nil
	}

	return &foundOrganizer
}

func (r *OrganizerMongoDBRepository) FindOrganizerByUsername(username string) *models.Organizer {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"username": username})

	var foundOrganizer models.Organizer
	err := result.Decode(&foundOrganizer)

	if err != nil {
		return nil
	}

	return &foundOrganizer
}

func (r *OrganizerMongoDBRepository) AnyWithID(organizerID uuid.UUID) bool {
	return r.FindOrganizer(organizerID) != nil
}

func (r *OrganizerMongoDBRepository) FindAll() []*models.Organizer {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result, err := organizers.Find(findContext, bson.M{})
	if err != nil {
		return nil
	}

	var foundOrganizers []*models.Organizer
	if err := result.Decode(&foundOrganizers); err != nil {
		return nil
	}

	return foundOrganizers
}
