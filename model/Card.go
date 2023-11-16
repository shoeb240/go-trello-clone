package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ListID      primitive.ObjectID `bson:"list_id" json:"list_id" validate:"required"`
	BoardID     primitive.ObjectID `bson:"board_id" json:"board_id" validate:"required"`
	UserID      string             `bson:"user_id" json:"user_id" validate:"required"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description"`
	Position    int                `bson:"position" json:"position"`
}
