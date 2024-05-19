package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Community struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	Description  string             `json:"description" bson:"description"`
	AvatarUrl    string             `bson:"avatar_url" bson:"avatar_url"`
	MembersCount int                `bson:"members_count"`
	// MembershipType: is it public, restricted, byInvite, etc

	CreatedByUserId primitive.ObjectID `bson:"created_by_user_id" bson:"created_by_user_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	IsDeleted       bool               `json:"is_deleted"`
}

// a user creates many communities. We're more interested in who created a given community than how many communities a user created
// user can create many posts
// a post belong to a community
