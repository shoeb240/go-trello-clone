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

func (r *ListRepository) FindByID(ctx context.Context, ID string) (*model.List, error) {
	var listModel model.List
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&listModel); err != nil {
		return nil, errors.New("this is error")
	}

	return &listModel, nil
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
		return "x", err
	}

	listID := ctx.Value(listIDKey).(string)
	listObjID, err := primitive.ObjectIDFromHex(listID)
	if err != nil {
		return "y", err
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
		return "Updated", nil
	}
	return "Not Modified", nil
}
