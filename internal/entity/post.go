package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Votes struct {
	Up   int `bson:"up""`
	Down int `bson:"down"`
}

type Comment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Content         string             `bson:"content"`
	CreatedByUserId primitive.ObjectID `bson:"created_by_user_id" bson:"created_by_user_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	IsDeleted       bool               `json:"is_deleted"`

	Votes Votes `bson:"votes"`
}

type Post struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`

	CreatedByUserId primitive.ObjectID `bson:"created_by_user_id" bson:"created_by_user_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	IsDeleted       bool               `json:"is_deleted"`
	Votes           Votes              `bson:"votes"`
	Comments        []Comment          `bson:"comments"`
}
