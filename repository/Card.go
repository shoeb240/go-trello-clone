package repository

import (
	"context"
	"errors"

	"github.com/shoeb240/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardRepository struct {
	collection *mongo.Collection
}

func newCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{
		collection: db.Collection("Card"),
	}
}

func (r *CardRepository) FindByID(ctx context.Context, ID string) (*model.Card, error) {
	var cardModel model.Card
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&cardModel); err != nil {
		return nil, errors.New("this is error")
	}

	return &cardModel, nil
}

func (r *CardRepository) Create(ctx context.Context, cardModel model.Card) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, cardModel)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *CardRepository) Update(ctx context.Context, cardModel model.Card) (string, error) {
	boardID := ctx.Value(boardIDKey).(string)
	objID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return "x", err
	}

	cardID := ctx.Value(listIDKey).(string)
	listObjID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return "y", err
	}

	filter := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": bson.M{
			"lists.$[elem].title":    cardModel.Title,
			"lists.$[elem].position": cardModel.Position,
		},
	}
	arrayFilters := options.Update().SetArrayFilters(
		options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem._id": listObjID}},
		},
	)
	result, err := r.collection.UpdateOne(ctx,
		filter,
		update,
		arrayFilters,
	)

	if err != nil {
		return "", err
	}
	if result.ModifiedCount > 0 {
		return "Updated", nil
	}
	return "Not Modified", nil
}
