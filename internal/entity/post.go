package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Votes struct {
	Up   int `bson:"up""`
	Down int `bson:"down"`
}

type Stats struct {
	Upvotes      int `bson:"upvotes" json:"upvotes"`
	Downvotes    int `bson:"downvotes" json:"downvotes"`
	CommentCount int `bson:"commentCount" json:"commentCount"`
	Awards       []struct {
		Type  string `bson:"type" json:"type"`
		Count int    `bson:"count" json:"count"`
	} `bson:"awards" json:"awards"`
}

type Metadata struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	Tags      []string  `bson:"tags" json:"tags"`
}

type Content struct {
	Type  string   `bson:"type" json:"type"` // text, link, image, video, poll
	Body  string   `bson:"body" json:"body"`
	Media []string `bson:"mediaUrls" json:"media"`
	Poll  struct {
		Options []string `bson:"options" json:"options"`
		Votes   []int    `bson:"votes" json:"votes"`
	} `bson:"poll,omitempty" json:"poll,omitempty"`
}

type Author struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}

type Post struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title         string             `bson:"title,omitempty" validate:"required"`
	Content       string             `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	CommunityID   primitive.ObjectID `bson:"communityId" json:"communityId"`
	CommunityName string             `bson:"communityName" json:"communityName"`
	Author        Author             `bson:"author" json:"author"`
	Metadata      Metadata           `bson:"metadata" json:"metadata"`
	MainContent   Content            `bson:"content" json:"content"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"` // Todo: deprecate
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"` // Todo: deprecate

	Stats           Stats              `bson:"stats" json:"stats"`
	CreatedByUserId primitive.ObjectID `json:"created_by_user_id" bson:"created_by_user_id"`

	IsDeleted bool  `json:"is_deleted" bson:"is_deleted"`
	Votes     Votes `bson:"votes"`

	//Flags struct {
	//	IsNSFW   bool `bson:"isNSFW" json:"isNSFW"`
	//	IsSpam   bool `bson:"isSpam" json:"isSpam"`
	//	IsPinned bool `bson:"isPinned" json:"isPinned"`
	//} `bson:"flags" json:"flags"`
}
