package mongoDBRepos

import (
	"ticken-event-service/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *EventMongoDBRepository) UpdateEvent(event *models.Event) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	updateOptions := new(options.FindOneAndUpdateOptions)
	updateOptions.SetReturnDocument(options.After)

	events := r.getCollection()
	result := events.FindOneAndUpdate(
		findContext,
		bson.M{"event_id": event.EventID},
		bson.M{
			"$set": bson.M{
				"name":            event.Name,
				"date":            event.Date.Format(time.RFC3339),
				"organization_id": event.OrganizationID,
				"pvt_bc_channel":  event.PvtBCChannel,
				"sections":        event.Sections,
				"on_chain":        event.OnChain,
				"on_sale":         event.OnSale,
				"status":          event.Status,
				"pub_bc_address":  event.PubBCAddress,
			},
		},
		updateOptions,
	)

	updatedEvent := new(models.Event)
	err := result.Decode(updatedEvent)
	if err != nil {
		return nil
	}
	return updatedEvent
}

// FindAvailableEvents
// Find all events that are on sale
func (r *EventMongoDBRepository) FindAvailableEvents() []*models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result, err := events.Find(findContext, bson.M{"status": bson.M{"$ne": models.EventStatusDraft}})

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
