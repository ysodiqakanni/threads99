package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TimelinePost struct {
	ID            string `bson:"_id"`
	Title         string `bson:"title"`
	CommunityId   string `bson:"communityId"`
	CommunityName string `bson:"communityName"`
	Content       struct {
		Body      string   `bson:"body"`
		Type      string   `bson:"type"`
		MediaUrls []string `bson:"mediaUrls"`
	} `bson:"content"`
	Author struct {
		ID       string `bson:"_id"`
		Username string
	} `bson:"author"`
	Metadata struct {
		CreatedAt time.Time `bson:"createdAt"`
	} `bson:"metadata"`
	Stats struct {
		Upvotes      int `bson:"upvotes" json:"upvotes"`
		Downvotes    int `bson:"downvotes" json:"downvotes"`
		CommentCount int `bson:"commentCount" json:"commentCount"`
	}
}

type PostResponseFlat struct {
	PostID        primitive.ObjectID `bson:"_id"`
	CommunityID   primitive.ObjectID `bson:"communityId"`
	CommunityName string             `bson:"communityName"`
	Title         string             `bson:"title"`
	ContentBody   string             `bson:"content.body"`
	ContentType   string             `bson:"content.type"`
	AuthorName    string             `bson:"author.username"`
	AuthorID      primitive.ObjectID `bson:"author._id"`
	TimePosted    time.Time          `bson:"created_at"`
	CommentCount  int32              `bson:"stats.commentCount"`
	VoteCount     int32              `bson:"-"` // Calculated field
}

type PostResponse struct {
	ID            primitive.ObjectID `bson:"_id"`
	Title         string             `bson:"title"`
	CommunityID   primitive.ObjectID `bson:"communityId"`
	CommunityName string             `bson:"communityName"`
	Author        Author             `bson:"author"`
	Content       Content            `bson:"content"`
	CreatedAt     time.Time          `bson:"created_at"`
	Stats         Stats              `bson:"stats"`
}

type Author struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
}

type Content struct {
	Type string `bson:"type"`
	Body string `bson:"body"`
}

type Stats struct {
	CommentCount int32 `bson:"commentCount"`
	UpVotes      int32 `bson:"upvotes"`
	DownVotes    int32 `bson:"downvotes"`
}

type PostLite struct {
}
