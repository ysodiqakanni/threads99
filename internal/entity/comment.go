package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentMetadata struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	//EditHistory []struct {
	//	Content   string    `bson:"content" json:"content"`
	//	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	//} `bson:"editHistory" json:"editHistory"`
}

// Todo: Do we need to include communityId and Name?
type Comment struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	PostID   primitive.ObjectID  `bson:"postId" json:"postId"`
	ParentID *primitive.ObjectID `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Author   Author              `bson:"author" json:"author"`
	Content  Content             `bson:"content" json:"content"`
	Metadata CommentMetadata     `bson:"metadata" json:"metadata"`
	Stats    Stats               `bson:"stats" json:"stats"`
	//Flags struct {
	//	IsDeleted bool `bson:"isDeleted" json:"isDeleted"`
	//	IsEdited  bool `bson:"isEdited" json:"isEdited"`
	//} `bson:"flags" json:"flags"`
}
