package repository

import (
	"context"
	"errors"
	"log"

	"github.com/shoeb240/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BoardRepository struct {
	collection *mongo.Collection
}

func newBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{
		collection: db.Collection("Board"),
	}
}

func (r *BoardRepository) FindByID(ctx context.Context, objID primitive.ObjectID) (model.Board, error) {
	var boardModel model.Board

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"_id": objID,
			},
		},
		bson.M{
			"$unwind": "$lists",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "Card",
				"localField":   "lists.cards",
				"foreignField": "_id",
				"as":           "lists.card_details",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":     "$_id",
				"user_id": bson.M{"$first": "$user_id"},
				"title":   bson.M{"$first": "$title"},
				"lists":   bson.M{"$push": "$lists"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&boardModel); err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return boardModel, nil
}

func (r *BoardRepository) Create(ctx context.Context, boardModel model.Board) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, boardModel)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *BoardRepository) AddCardToList(ctx context.Context, cardModel model.Card) error {
	filter := bson.M{
		"_id": cardModel.BoardID,
	}
	update := bson.M{
		"$push": bson.M{
			"lists.$[elem].cards": cardModel.ID,
		},
	}
	arrayFilters := options.Update().SetArrayFilters(
		options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem._id": cardModel.ListID}},
		},
	)
	result, err := r.collection.UpdateOne(ctx,
		filter,
		update,
		arrayFilters,
	)

	if err != nil {
		return err
	}
	if result.ModifiedCount <= 0 {
		return errors.New("list not updated")
	}
	return nil
}

func (r *BoardRepository) Update(ctx context.Context, updateData primitive.M) (string, error) {
	boardID := ctx.Value(boardIDKey).(string)
	boardObjID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": boardObjID,
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

func (r *BoardRepository) RemoveCardFromList(ctx context.Context, cardDeleteFromList CardDeleteFromList) error {
	filter := bson.M{
		"_id": cardDeleteFromList.BoardID,
	}
	update := bson.M{
		"$pull": bson.M{
			"lists.$[elem].cards": cardDeleteFromList.CardID,
		},
	}
	arrayFilters := options.Update().SetArrayFilters(
		options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem._id": cardDeleteFromList.ListID}},
		},
	)
	result, err := r.collection.UpdateOne(ctx,
		filter,
		update,
		arrayFilters,
	)

	if err != nil {
		return err
	}
	if result.ModifiedCount <= 0 {
		return errors.New("card is not removed from list")
	}
	return nil
}

func (r *BoardRepository) Delete(ctx context.Context, boardObjID primitive.ObjectID) (string, error) {
	filter := bson.M{
		"_id": boardObjID,
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
