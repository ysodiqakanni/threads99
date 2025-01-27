package dto

import "time"

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
		ID string `bson:"_id"`
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
