package user

import (
	"context"
	"fmt"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

// Repository encapsulates the logic to access categories from the data source.
type Repository interface {
	Get(ctx context.Context, id primitive.ObjectID) (entity.User, error)
	GetByEmail(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user entity.User) (*primitive.ObjectID, error)
	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
	IsUserExistsByUsername(ctx context.Context, email string) (bool, error)
	StartSession() (mongo.Session, error)
}

// repository persists albums in database
type repository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	col := db.DB().Collection("users")
	return repository{col, logger}
}

func (r repository) StartSession() (mongo.Session, error) {
	return r.collection.Database().Client().StartSession()
}

func (r repository) Get(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	filter := bson.M{"_id": id}
	var user entity.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}
func (r repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	filter := bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "^" + email + "$", Options: "i"}}}
	var user entity.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	// any other error apart from the ErrNoDocuments is something to be worried about.

	//fmt.Println("user data: ", user)
	return &user, err
}
func (r repository) Create(ctx context.Context, user entity.User) (*primitive.ObjectID, error) {
	// save user email to lowercase to avoid extra conversion during lookup
	user.Email = strings.ToLower(user.Email)
	user.Username = strings.ToLower(user.Username)
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	fmt.Printf("inserted user data with ID %v\n", result.InsertedID)
	id := result.InsertedID.(primitive.ObjectID)
	return &id, err
}

func (r repository) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	// Note that emails are saved in lowercase.
	filter := bson.M{"email": strings.ToLower(email)}
	opts := options.Count().SetLimit(1)
	count, err := r.collection.CountDocuments(ctx, filter, opts)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r repository) IsUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	// Note that emails are saved in lowercase.
	filter := bson.M{"username": strings.ToLower(username)}
	opts := options.Count().SetLimit(1)
	count, err := r.collection.CountDocuments(ctx, filter, opts)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
