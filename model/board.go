package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"user_id" validate:"required"`
	Title  string             `bson:"title" json:"title" validate:"required"`
	Lists  []List             `bson:"lists" json:"lists" validate:"required"`
}
