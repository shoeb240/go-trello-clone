package repository

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func newUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("User"),
	}
}

func (r *UserRepository) Signup(ctx context.Context, userModel model.User) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, userModel)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
