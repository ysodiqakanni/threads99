package post

import (
	"context"
	"errors"
	"fmt"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id primitive.ObjectID) (entity.Post, error)
	Create(ctx context.Context, postRequest entity.Post) (*primitive.ObjectID, error)
	AddCommentToPost(ctx context.Context, postId primitive.ObjectID, comment entity.Comment) error
	UpvoteComment(ctx context.Context, commentId primitive.ObjectID, postId primitive.ObjectID, voteValue int) error
	UpvotePost(ctx context.Context, postId primitive.ObjectID, voteValue int) error
	GetCommentsByPostId(ctx context.Context, postId primitive.ObjectID) ([]entity.Comment, error)
}

// repository persists data in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("posts")
	logger.Infof("collection retrieved")
	return repository{col, logger}
}

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (entity.Post, error) {
	fmt.Println("Getting post by Id")
	filter := bson.M{"_id": id}
	var post entity.Post
	err := r.collection.FindOne(ctx, filter).Decode(&post)

	return post, err
}

func (r repository) Create(ctx context.Context, postRequest entity.Post) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, postRequest)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}

func (r repository) AddCommentToPost(ctx context.Context, postId primitive.ObjectID, comment entity.Comment) error {
	filter := bson.M{"_id": postId}
	post, err := r.Get(ctx, postId)
	if err != nil {
		// maybe the post doesn't exist?
		return err
	}
	if post.Comments == nil {
		initCommentUpdate := bson.M{
			"$set": bson.M{"comments": []entity.Comment{}},
		}
		_, err = r.collection.UpdateOne(ctx, filter, initCommentUpdate)
		if err != nil {
			return nil
		}
	}
	update := bson.M{
		"$push": bson.M{"comments": comment},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r repository) UpvoteComment(ctx context.Context, commentId primitive.ObjectID, postId primitive.ObjectID, voteValue int) error {

	filter := bson.M{
		"_id":          postId,
		"comments._id": commentId,
	}
	update := bson.M{
		"$inc": bson.M{"comments.$.votes.up": voteValue},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r repository) UpvotePost(ctx context.Context, postId primitive.ObjectID, voteValue int) error {

	filter := bson.M{
		"_id": postId,
	}
	update := bson.M{
		"$inc": bson.M{"votes.up": voteValue},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r repository) GetCommentsByPostId(ctx context.Context, postId primitive.ObjectID) ([]entity.Comment, error) {
	filter := bson.M{
		"_id": postId,
	}
	var result struct {
		Comments []entity.Comment `bson:"comments"`
	}
	err := r.collection.FindOne(ctx, filter).Decode(&result)

	return result.Comments, err

	var post bson.M
	err = r.collection.FindOne(ctx, filter).Decode(&post)

	emptyComments := []entity.Comment{}
	if err != nil {
		return emptyComments, err
	}
	comments, ok := post["comments"].([]interface{})
	if !ok {
		r.logger.Errorf("comments not found or invalid")
		return emptyComments, errors.New("An error occurred while fetching comments from doc.")
	}

	var commentsList []entity.Comment
	for _, comment := range comments {
		var commentDoc entity.Comment
		bsonBytes, err := bson.Marshal(comment)
		if err != nil {
			continue
		}
		err = bson.Unmarshal(bsonBytes, &commentDoc)
		if err != nil {
			continue
		}
		commentsList = append(commentsList, commentDoc)
	}
	return commentsList, nil
}
