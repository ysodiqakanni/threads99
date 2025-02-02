package comment

import (
	"context"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateNewComment(ctx context.Context, comment entity.Comment) (error, *primitive.ObjectID)
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
