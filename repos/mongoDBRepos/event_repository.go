package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
	"time"
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
			primKeyField:   "event_id",
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

func (r *EventMongoDBRepository) FindEvent(eventID uuid.UUID) *models.Event {
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

func (r *EventMongoDBRepository) FindOrganizationEvents(organizationID uuid.UUID) []*models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result, err := events.Find(findContext, bson.M{"organization_id": organizationID})

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

func (r *EventMongoDBRepository) UpdateEventStatus(event *models.Event) error {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	updateOptions := new(options.FindOneAndUpdateOptions)
	updateOptions.SetReturnDocument(options.After)

	events := r.getCollection()
	result := events.FindOneAndUpdate(
		findContext,
		bson.M{"event_id": event.EventID},
		bson.M{"$set": bson.M{"status": event.Status}},
		updateOptions,
	)

	updatedEvent := new(models.Event)
	err := result.Decode(updatedEvent)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventMongoDBRepository) AddSectionToEvent(eventID uuid.UUID, section *models.Section) error {
	updateContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	filter := bson.M{"event_id": eventID}
	update := bson.M{"$push": bson.M{"sections": section}}

	_, err := events.UpdateOne(updateContext, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventMongoDBRepository) UpdatePUBBCData(event *models.Event) error {
	updateContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	filter := bson.M{"event_id": event.EventID}
	update := bson.M{"$set": bson.M{
		"pubbc_address": event.PubBCAddress,
		"pubbc_tx_id":   event.PubBCTxID,
	}}

	_, err := events.UpdateOne(updateContext, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventMongoDBRepository) FindEvents(withName string, withStatus []models.EventStatus, fromDate time.Time, toDate time.Time) []*models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()

	filter := bson.M{}

	if len(withName) > 0 {
		filter["name"] = "/a/"
	}
	if len(withStatus) > 0 {
		filter["status"] = bson.M{"$in": withStatus}
	}
	if !fromDate.IsZero() {
		filter["date"] = bson.M{"$ge": fromDate}
	}
	if !toDate.IsZero() {
		filter["date"] = bson.M{"$le": toDate}
	}

	result, err := events.Find(findContext, filter)
	if err != nil {
		return nil
	}

	var foundEvents = make([]*models.Event, 0)
	for result.Next(findContext) {
		var event = new(models.Event)
		err = result.Decode(event)
		foundEvents = append(foundEvents, event)
	}

	return foundEvents
}
