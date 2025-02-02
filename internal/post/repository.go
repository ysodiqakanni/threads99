package post

import (
	"context"
	"errors"
	"fmt"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Get(ctx context.Context, id primitive.ObjectID) (entity.Post, error)
	Create(ctx context.Context, postRequest entity.Post) (*primitive.ObjectID, error)
	//AddCommentToPost(ctx context.Context, postId primitive.ObjectID, comment entity.Comment) error
	UpvoteComment(ctx context.Context, commentId primitive.ObjectID, postId primitive.ObjectID, voteValue int) error
	UpvotePost(ctx context.Context, postId primitive.ObjectID, voteValue int) error
	GetCommentsByPostId(ctx context.Context, postId primitive.ObjectID) ([]entity.Comment, error)
	GetAllRecentPosts(ctx context.Context) ([]dto.TimelinePost, error)
	GetPostLiteById(ctx context.Context, id primitive.ObjectID) (*dto.PostResponse, error)
}

// repository persists data in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

// Called in main.go to initialize the repo.
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

func (r repository) GetPostLiteById(ctx context.Context, id primitive.ObjectID) (*dto.PostResponse, error) {
	filter := bson.M{"_id": id}

	// Define projection to fetch only needed fields
	projection := bson.M{
		"_id":                1,
		"title":              1,
		"communityId":        1,
		"communityName":      1,
		"author._id":         1,
		"author.username":    1,
		"content.type":       1,
		"content.body":       1,
		"created_at":         1,
		"stats.commentCount": 1,
		"stats.upvotes":      1,
		"stats.downvotes":    1,
	}

	opts := options.FindOne().SetProjection(projection)

	var post dto.PostResponse
	err := r.collection.FindOne(ctx, filter, opts).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil

	///
	/*

		projection := bson.D{
			{"_id", 1},
			{"communityId", 1},
			{"communityName", 1},
			{"title", 1},
			{"content.body", 1},
			{"content.type", 1},
			{"author.username", 1},
			{"author._id", 1},
			{"created_at", 1},
			{"stats.commentCount", 1},
			{"votes.up", 1},
			{"votes.down", 1},
		}

		// Create filter for finding by postID
		filter := bson.M{"_id": id}

		// Execute the query
		var post dto.PostResponse
		err := r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&post)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}

		return &post, nil
	*/
}

// / SECTION POSTS
func (r repository) Create(ctx context.Context, postRequest entity.Post) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, postRequest)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}

// Todo: Should be removed! Never get ALL!!!
func (r repository) GetAllRecentPosts(ctx context.Context) ([]dto.TimelinePost, error) {
	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "communityId", Value: 1},
		{Key: "communityName", Value: 1},
		{Key: "content.body", Value: 1},
		{Key: "content.type", Value: 1},
		{Key: "content.mediaUrls", Value: 1},
		{Key: "author._id", Value: 1},
		{Key: "metadata.createdAt", Value: 1},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetProjection(projection).
		SetLimit(30)

	// Execute the query
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// Decode the results into our lightweight struct
	var posts []dto.TimelinePost
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
	// Previous codes
	/*
		// Create options to sort by creation time in descending order (-1)
		opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
		// Find all documents with sorting options
		cursor, err := r.collection.Find(ctx, bson.M{}, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		// Decode all documents into a slice of posts
		var posts []entity.Post
		if err = cursor.All(ctx, &posts); err != nil {
			return nil, err
		}

		return posts, nil
	*/
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

// SECTION POSTS
/*
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
*/

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
