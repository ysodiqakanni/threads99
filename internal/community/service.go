package community

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"time"
)

type Community struct {
	entity.Community
}
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (*Community, error)
	Create(ctx context.Context, req CreateCommunityRequest) (*Community, error)
	GetAllCommunities(ctx context.Context) ([]entity.Community, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new post service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type CreateCommunityRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	AvatarUrl       string `bson:"avatar_url" bson:"avatar_url"`
	CreatedByUserId string
}

func (m CreateCommunityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128), validation.Match(regexp.MustCompile("^[a-zA-Z0-9].*$"))),
	)
}

func (s service) Create(ctx context.Context, req CreateCommunityRequest) (*Community, error) {
	now := time.Now()
	userId, err := primitive.ObjectIDFromHex(req.CreatedByUserId)
	if err != nil {
		// invalid userId
		return nil, err
	}
	// check if the community already exists
	existingCommunity, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		// db error
		// Todo: log it
		s.logger.Error(err)
		return nil, err
	}
	if existingCommunity.Name != "" {
		// already exists
		return nil, errors.New("A community with this name already exists.")
	}
	id, err := s.repo.Create(ctx, entity.Community{
		Name:            req.Name,
		Description:     req.Description,
		CreatedByUserId: userId,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, *id)
}

func (s service) Get(ctx context.Context, id primitive.ObjectID) (*Community, error) {
	community, err := s.repo.Get(ctx, id)
	if err != nil || community == nil {
		// check if it's
		return nil, err
	}
	return &Community{*community}, nil
}

func (s service) GetAllCommunities(ctx context.Context) ([]entity.Community, error) {
	communities, err := s.repo.GetAllCommunities(ctx)
	if err != nil {
		return nil, err
	}

	return communities, nil
}
