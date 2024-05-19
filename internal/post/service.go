package post

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

type CreateNewPostRequest struct {
	Content         string `json:"content"`
	CreatedByUserId string `json:"created_by_user_id"`
}

func (m CreateNewPostRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Content, validation.Required, validation.Length(0, 1024)),
		validation.Field(&m.CreatedByUserId, validation.Required),
	)
}

// Get returns the post with the specified the post ID.
func (s service) Get(ctx context.Context, id primitive.ObjectID) (Post, error) {
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return Post{}, err
	}
	return Post{post}, nil
}

func (s service) CreatePost(ctx context.Context, request CreateNewPostRequest) error {
	post := entity.Post{
		Content: request.Content,
		//CreatedByUserId: request.CreatedByUserId
	}

	rez, err := s.repo.Create(ctx, post)
	if err != nil {

	}
}
