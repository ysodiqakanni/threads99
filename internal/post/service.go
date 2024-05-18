package post

import (
	"context"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	entity.Post
}

type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (Post, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new post service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the post with the specified the post ID.
func (s service) Get(ctx context.Context, id primitive.ObjectID) (Post, error) {
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return Post{}, err
	}
	return Post{post}, nil
}
