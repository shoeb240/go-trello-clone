package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title" validate:"required"`
	Position int                `bson:"position" json:"position"`
	Cards    []Card             `bson:"cards" json:"cards" validate:"required"`
}
