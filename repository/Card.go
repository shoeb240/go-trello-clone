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

type CardMoveReq struct {
	CardID     primitive.ObjectID `bson:"card_id" json:"card_id"`
	BoardID    primitive.ObjectID `bson:"board_id" json:"board_id" validate:"required"`
	FromListID primitive.ObjectID `bson:"from_list_id" json:"from_list_id" validate:"required"`
	ToListID   primitive.ObjectID `bson:"to_list_id" json:"to_list_id" validate:"required"`
	ToPosition int                `bson:"to_position" json:"to_position" validate:"required"`
}

type CardDeleteFromList struct {
	CardID  primitive.ObjectID
	BoardID primitive.ObjectID
	ListID  primitive.ObjectID
}

func newCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{
		collection: db.Collection("Card"),
	}
}

func (r *CardRepository) FindByID(ctx context.Context, objID primitive.ObjectID) (model.Card, error) {
	var cardModel model.Card

	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&cardModel); err != nil {
		return cardModel, errors.New("card not found")
	}

	return cardModel, nil
}

func (r *CardRepository) Create(ctx context.Context, cardModel model.Card) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, cardModel)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *CardRepository) Update(ctx context.Context, updateData primitive.M) (string, error) {
	cardID := ctx.Value(cardIDKey).(string)
	cardObjID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": cardObjID,
	}
	update := bson.M{
		"$set": updateData,
	}
	result, err := r.collection.UpdateOne(ctx,
		filter,
		update,
	)

	if err != nil {
		return "", err
	}
	if result.ModifiedCount > 0 {
		return "updated", nil
	}
	return "not modified", nil
}

func (r *CardRepository) Delete(ctx context.Context, cardObjID primitive.ObjectID) (string, error) {
	filter := bson.M{
		"_id": cardObjID,
	}
	result, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return "", err
	}
	if result.DeletedCount > 0 {
		return "deleted", nil
	}
	return "not deleted", nil
}
