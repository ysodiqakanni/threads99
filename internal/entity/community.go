package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Community struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"name" bson:"name"`
	//CreatedByUserId
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}
