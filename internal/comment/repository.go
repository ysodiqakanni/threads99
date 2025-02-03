package comment

import (
	"context"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
)

type Repository interface {
	CreateNewComment(ctx context.Context, comment entity.Comment) (error, *primitive.ObjectID)
	GetCommentsByPostId(ctx context.Context, postId primitive.ObjectID) ([]dto.CommentTree, error)
}

// repository persists data in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("comments")
	logger.Infof("collection retrieved")
	return repository{col, logger}
}

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}

func (r repository) CreateNewComment(ctx context.Context, comment entity.Comment) (error, *primitive.ObjectID) {
	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return err, nil
	}
	id := result.InsertedID.(primitive.ObjectID)
	return err, &id
}
func (r repository) GetCommentsByPostId(ctx context.Context, postId primitive.ObjectID) ([]dto.CommentTree, error) {
	// Todo: throw error if the post does not exist.
	// Todo: add pagination for comments.
	var comments []entity.Comment
	cursor, err := r.collection.Find(ctx, bson.M{
		"postId": postId,
	})

	// the parent comments to return. Each can have one or more children.
	var parentComments []dto.CommentTree
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return parentComments, nil
		}
		return nil, err
	}

	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	// now add each parentComment (with their children) to the result array.
	// loop through the comments.
	// if current has a parent, add it to the parent's children.
	// else, make it a new parent with empty children.

	// since we need a list of comment->Comments
	// the result contains only comments with no parentId.
	// other comments are added to their parent using the parentIds on them.

	// Now let's create a map of id to comment object.
	// now we loop, if you're a child parent, you're added to the nodes of your parent.
	// otherwise, you're added to the result.

	commentMap := make(map[primitive.ObjectID]dto.CommentTree) // we could just use a list but dict is good for O(1) lookup
	for _, comment := range comments {
		commentTree := dto.CommentTree{
			Comment: comment,
			Replies: []dto.CommentTree{},
		}
		commentMap[comment.ID] = commentTree
	}

	for _, comment := range comments {
		if comment.ParentID != nil {
			// This is a nested reply. Find its parent and add to its replies.
			if _, exists := commentMap[*comment.ParentID]; exists {
				parentTree := commentMap[*comment.ParentID]                             // the parent comment.
				parentTree.Replies = append(parentTree.Replies, commentMap[comment.ID]) // the replies being appended to
				commentMap[*comment.ParentID] = parentTree                              // save back the parent.
			}
		} else {
			// A root or parent comment found! Add to the results.
			//parentComments = append(parentComments, commentMap[comment.ID])
		}
	}

	// Third pass:
	// Todo: fix pointer issue in prior loop to avoid this extra loop.
	for _, comment := range comments {
		if comment.ParentID == nil {
			parentComments = append(parentComments, commentMap[comment.ID])
		}
	}

	sort.Slice(parentComments, func(i, j int) bool {
		return parentComments[i].Comment.Metadata.CreatedAt.After(
			parentComments[j].Comment.Metadata.CreatedAt)
	})

	return parentComments, nil
}
