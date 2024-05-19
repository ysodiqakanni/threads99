package user

import (
	"context"
	"errors"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service encapsulates use case logic for businessCategories.
type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
}

// User represents the data about a User.
type User struct {
	entity.User
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new category service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context, id primitive.ObjectID) (*User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}

func (s service) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	return User{user}, err

	//if err != nil {
	//	// somt
	//	return User{}, err
	//}
	//return User{user}, nil
}
func (s service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, getErr := s.GetByEmail(ctx, req.Email)

	if getErr != nil {
		return nil, errors.New("An unknown error occurred while fetching user data")
	}

	if existing.Email != "" {
		// an empty object should return an empty email (and other props)
		return nil, errors.New("A user with this name already exists")
	}

	//password :=
	// Todo: generate user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id, err := s.repo.Create(ctx, entity.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		Role:           req.Roles,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, *id)
}
