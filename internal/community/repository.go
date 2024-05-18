package community

import (
	"context"
	"fmt"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, community entity.Community) (*primitive.ObjectID, error)
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
