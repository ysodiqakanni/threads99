package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`

	Email          string   `bson:"email"`
	Username       string   `bson:"username"`
	Role           []string `bson:"role"`
	HashedPassword []byte   `bson:"hashed_password"`

	Bio           string
	CoverPhotoUrl string
	LogoUrl       string
	PostsCount    int
	CommentsCount int

	JoinedCommunities   []string `bson:"joined_communities"` // array of IDs of communities joined by user
	FavoriteCommunities []string `bson:"favorite_communities"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

// Identity represents an authenticated user identity.
type UserAuthIdentity interface {
	// GetID returns the user ID.
	GetID() primitive.ObjectID
	GetUserName() string
	GetEmail() string
	GetRole() []string
}

// implement Identity's functions for auth purpose
func (u User) GetRole() []string {
	return u.Role
}

// GetID returns the user ID.
func (u User) GetID() primitive.ObjectID {
	return u.ID
}

// GetName returns the user name.
func (u User) GetUserName() string {
	return u.Username
}
func (u User) GetEmail() string {
	return u.Email
}
