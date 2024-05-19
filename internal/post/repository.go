package post

import (
	"context"
	"fmt"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Get(ctx context.Context, id primitive.ObjectID) (entity.Post, error)
	Create(ctx context.Context, postRequest entity.Post) (*primitive.ObjectID, error)
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
