package mongoDBRepos

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
)

const EventCollectionName = "events"

type EventMongoDBRepository struct {
	baseRepository
}

func NewEventRepository(dbClient *mongo.Client, dbName string) *EventMongoDBRepository {

	return &EventMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: EventCollectionName,
		},
	}
}

func (r *EventMongoDBRepository) getCollection() *mongo.Collection {
	ctx, cancel := r.generateOpSubcontext()
	defer cancel()

	coll := r.baseRepository.getCollection()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "event_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("error creating event repo: " + err.Error())
	}
	return coll
}

func (r *EventMongoDBRepository) AddEvent(event *models.Event) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	_, err := events.InsertOne(storeContext, event)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventMongoDBRepository) FindEvent(eventID string) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"event_id": eventID})

	var foundEvent models.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}

func (r *EventMongoDBRepository) FindOrgEvents(orgID string) []*models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result, err := events.Find(findContext, bson.M{"organizer_id": orgID})

	if err != nil {
		return nil
	}

	var foundEvents []*models.Event
	for result.Next(findContext) {
		var event = new(models.Event)
		err = result.Decode(event)
		foundEvents = append(foundEvents, event)
	}
	return foundEvents
}

func (r *EventMongoDBRepository) UpdateEvent(EventID string, OrganizerID string, PvtBCChannel string, Sections []models.Section) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	updateOptions := new(options.FindOneAndUpdateOptions)
	updateOptions.SetReturnDocument(options.After)

	fmt.Printf("event sections %s", Sections)

	events := r.getCollection()
	result := events.FindOneAndUpdate(
		findContext,
		bson.M{"event_id": EventID},
		bson.M{
			"$set": bson.M{
				"organizer_id":   OrganizerID,
				"pvt_bc_channel": PvtBCChannel,
				"sections":       Sections,
			},
		},
		updateOptions,
	)

	updatedEvent := new(models.Event)
	err := result.Decode(updatedEvent)
	if err == nil {
		return nil
	}
	return updatedEvent
}