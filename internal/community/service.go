package community

import (
	"context"
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
	Get(ctx context.Context, id primitive.ObjectID) (Community, error)
	Create(ctx context.Context, req CreateCommunityRequest) (Community, error)
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
}

func (m CreateCommunityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128), validation.Match(regexp.MustCompile("^[a-zA-Z0-9].*$"))),
	)
}

func (s service) Create(ctx context.Context, req CreateCommunityRequest) (Community, error) {
	//if err := req.Validate(); err != nil {
	//	return Community{}, err
	//}

	//existing, _ := s.GetByName(ctx, req.Name)
	////emptyObj := BusinessCategory{}
	//if existing != nil /*!= emptyObj*/ {
	//	return BusinessCategory{}, errors.New("A business_ category with this name already exists")
	//}

	now := time.Now()
	id, err := s.repo.Create(ctx, entity.Community{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return Community{}, err
	}
	return s.Get(ctx, *id)
}

func (s service) Get(ctx context.Context, id primitive.ObjectID) (Community, error) {
	post, err := s.repo.Get(ctx, id)
	if err != nil {
		return Community{}, err
	}
	return Community{post}, nil
}
