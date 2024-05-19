package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	FirstName           string             `bson:"first_name"`
	LastName            string             `bson:"last_name"`
	JoinedCommunities   []string           `bson:"joined_communities"` // array of IDs of communities joined by user
	FavoriteCommunities []string           `bson:"favorite_communities"`

	Email          string   `bson:"email"`
	Role           []string `bson:"role"`
	HashedPassword []byte   `bson:"hashed_password"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}
