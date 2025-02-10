package community

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
	Create(ctx context.Context, community entity.Community) (*primitive.ObjectID, error)
	Get(ctx context.Context, id primitive.ObjectID) (*entity.Community, error)
	GetByName(ctx context.Context, name string) (entity.Community, error)
	GetAllCommunities(ctx context.Context) ([]entity.Community, error)
}

// repository persists data in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("communities")
	logger.Infof("collection retrieved")
	return repository{col, logger}
}

func (r repository) Create(ctx context.Context, community entity.Community) (*primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, community)
	if err != nil {
		return nil, err
	}

	fmt.Printf("inserted document with ID %v\n", result.InsertedID)
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (*entity.Community, error) {
	fmt.Println("Getting community by Id")
	filter := bson.M{"_id": id}
	var community entity.Community
	err := r.collection.FindOne(ctx, filter).Decode(&community)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		fmt.Println("Error fetching community with id: " + id.Hex())
		return nil, err
	}

	return &community, nil
}

// implement getCommunityByName
func (r repository) GetByName(ctx context.Context, name string) (entity.Community, error) {
	fmt.Println("Getting community by name")
	filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^" + name + "$", Options: "i"}}}
	var community entity.Community
	err := r.collection.FindOne(ctx, filter).Decode(&community)
	if err == mongo.ErrNoDocuments {
		return entity.Community{}, nil
	}

	return community, err
}

// Todo: Should be removed! Never get ALL!!!
func (r repository) GetAllCommunities(ctx context.Context) ([]entity.Community, error) {
	// Find all documents in the collection
	cursor, err := r.collection.Find(ctx, bson.M{}) // using an empty filter to get ALL
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode all documents into a slice of communities
	var communities []entity.Community
	if err = cursor.All(ctx, &communities); err != nil {
		return nil, err
	}

	return communities, nil
}
