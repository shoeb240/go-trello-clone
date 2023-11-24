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

type ListRepository struct {
	collection *mongo.Collection
}

type CustomKeyType string

const (
	boardIDKey CustomKeyType = "boardID"
	listIDKey  CustomKeyType = "listID"
)

func newListRepository(db *mongo.Database) *ListRepository {
	return &ListRepository{
		collection: db.Collection("Board"),
	}
}

func (r *ListRepository) Create(ctx context.Context, listModel model.List) error {
	boardID := ctx.Value(boardIDKey).(string)
	objID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateByID(ctx, objID, bson.M{"$push": bson.M{"lists": listModel}})
	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 {
		return nil
	}
	return errors.New("Failed")
}

func (r *ListRepository) Update(ctx context.Context, listModel model.List) (string, error) {
	boardID := ctx.Value(boardIDKey).(string)
	objID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return "", err
	}

	listID := ctx.Value(listIDKey).(string)
	listObjID, err := primitive.ObjectIDFromHex(listID)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": bson.M{
			"lists.$[elem].title":    listModel.Title,
			"lists.$[elem].position": listModel.Position,
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
		return "updated", nil
	}
	return "not modified", nil
}

func (r *ListRepository) Delete(ctx context.Context) (string, error) {
	boardID := ctx.Value(boardIDKey).(string)
	objID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return "", err
	}

	listID := ctx.Value(listIDKey).(string)
	listObjID, err := primitive.ObjectIDFromHex(listID)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$unset": bson.M{
			"lists.$[elem]": 1,
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
		return "deleted", nil
	}
	return "not deleted", nil
}

func (r *ListRepository) HasCard(ctx context.Context) (bool, error) {
	boardID := ctx.Value(boardIDKey).(string)
	objID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return false, err
	}

	listID := ctx.Value(listIDKey).(string)
	listObjID, err := primitive.ObjectIDFromHex(listID)
	if err != nil {
		return false, err
	}

	var hasCard bool
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": objID},
		},
		bson.M{
			"$project": bson.M{"lists": bson.M{"$filter": bson.M{
				"input": "$lists",
				"as":    "list",
				"cond":  bson.M{"$eq": bson.A{"$$list._id", listObjID}},
			}}},
		},
		bson.M{
			"$unwind": "$lists",
		},
		bson.M{
			"$project": bson.M{
				"_id":      0,
				"numCards": bson.M{"$size": "$lists.cards"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return hasCard, err
	}
	defer cursor.Close(ctx)

	var result struct {
		NumCards int `bson:"numCards"`
	}

	cursor.Next(ctx)
	if err := cursor.Decode(&result); err != nil {
		return hasCard, err
	}

	return result.NumCards > 0, nil
}
