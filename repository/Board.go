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

type BoardRepository struct {
	collection *mongo.Collection
}

func newBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{
		collection: db.Collection("Board"),
	}
}

func (r *BoardRepository) FindByID(ctx context.Context, ID string) (model.Board, error) {
	var boardModel model.Board
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return model.Board{}, err
	}

	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&boardModel); err != nil {
		return model.Board{}, errors.New("this is error")
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
