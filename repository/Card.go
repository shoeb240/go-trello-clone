package repository

import (
	"context"
	"errors"

	"github.com/shoeb240/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardRepository struct {
	collection *mongo.Collection
}

const (
	cardIDKey CustomKeyType = "cardID"
)

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
	cardID := ctx.Value(cardIDKey).(string)
	cardObjID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": cardObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"title":       cardModel.Title,
			"description": cardModel.Description,
			"position":    cardModel.Position,
		},
	}
	result, err := r.collection.UpdateOne(ctx,
		filter,
		update,
	)

	if err != nil {
		return "", err
	}
	if result.ModifiedCount > 0 {
		return "Updated", nil
	}
	return "Not Modified", nil
}
