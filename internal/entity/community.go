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
	CoverPhotoUrl   string
	LogoUrl         string
	CreatedByUserId primitive.ObjectID `json:"CreatedByUserId" bson:"created_by_user_id" validate:"required"`
	CreatedAt       time.Time          `json:"CreatedAt"`
	UpdatedAt       time.Time          `json:"UpdatedAt"`
	IsDeleted       bool               `json:"IsDeleted"`
}

// a user creates many communities. We're more interested in who created a given community than how many communities a user created
// user can create many posts
// a post belong to a community
