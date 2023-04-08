package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-event-service/models"
	"ticken-event-service/utils"
)

const ValidatorCollectionName = "validators"

type ValidatorMongoDBRepository struct {
	baseRepository
}

func NewValidatorRepository(dbClient *mongo.Client, dbName string) *ValidatorMongoDBRepository {
	return &ValidatorMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: ValidatorCollectionName,
			primKeyField:   utils.GetStructTag(models.Validator{}.ValidatorID, "bson"),
		},
	}
}

func (r *ValidatorMongoDBRepository) getCollection() *mongo.Collection {
	ctx, cancel := r.generateOpSubcontext()
	defer cancel()

	coll := r.baseRepository.getCollection()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "validator_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("error creating event repo: " + err.Error())
	}
	return coll
}

func (r *ValidatorMongoDBRepository) AddValidator(validator *models.Validator) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	validators := r.getCollection()
	_, err := validators.InsertOne(storeContext, validator)
	if err != nil {
		return err
	}

	return nil
}

func (r *ValidatorMongoDBRepository) FindValidator(validatorID uuid.UUID) *models.Validator {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	organizers := r.getCollection()
	result := organizers.FindOne(findContext, bson.M{"validator_id": validatorID})

	var foundValidator models.Validator
	err := result.Decode(&foundValidator)

	if err != nil {
		return nil
	}

	return &foundValidator
}
