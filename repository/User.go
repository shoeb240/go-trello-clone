package repository

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	user := model.User{}

	filter := bson.M{
		"email": email,
	}

	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, userObjID primitive.ObjectID) (model.User, error) {
	user := model.User{}

	filter := bson.M{
		"_id": userObjID,
	}

	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}
