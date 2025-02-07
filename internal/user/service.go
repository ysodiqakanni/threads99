package user

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/ysodiqakanni/threads99/internal/dto"
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
	Create(ctx context.Context, req dto.CreateNewUserRequestDto) (*dto.CreateNewUserResponseDto, error)
	Login(ctx context.Context, username, password string) (string, error)
	GenerateJWT(identity entity.UserAuthIdentity) (string, error)
}

// User represents the data about a User.
type User struct {
	entity.User
}

type service struct {
	repo            Repository
	logger          log.Logger
	signingKey      string
	tokenExpiration int
}

// NewService creates a new category service.
func NewService(repo Repository, logger log.Logger, jwtSigningKey string, tokenExpiration int) Service {
	return service{repo, logger, jwtSigningKey, tokenExpiration}
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

	return User{*user}, err
}
func (s service) Create(ctx context.Context, req dto.CreateNewUserRequestDto) (*dto.CreateNewUserResponseDto, error) {
	if req.Username == "" {
		req.Username = req.Email
	}
	userExists, err := s.repo.IsUserExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("An unknown error occurred while fetching user data")
	}
	if userExists == true {
		// an empty object should return an empty email (and other props)
		return nil, errors.New("A user with this email already exists")
	}

	userExists, err = s.repo.IsUserExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("An unknown error occurred while fetching user data")
	}
	if userExists == true {
		// an empty object should return an empty email (and other props)
		return nil, errors.New("A user with this username already exists")
	}

	//password :=
	// Todo: generate user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	// Todo: design logic to generate username when it's not sent. Just extract part of the email and add some chars
	now := time.Now()
	id, err := s.repo.Create(ctx, entity.User{
		HashedPassword: hashedPassword,
		Email:          req.Email,
		Username:       req.Username,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return nil, err
	}

	ret := dto.CreateNewUserResponseDto{
		UserId:       id.Hex(),
		UserObjectId: *id,
		UserName:     req.Username,
	}
	return &ret, nil
}

func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if identity := s.authenticate(ctx, username, password); identity != nil {
		return s.GenerateJWT(identity)
	}
	return "", errors.New("Invalid username or password")
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, email, password string) entity.UserAuthIdentity {
	logger := s.logger.With(ctx, "user", email)

	// first get user by email
	// Todo: get by username as well.
	usr, err := s.repo.GetByEmail(ctx, email)
	if err != nil || usr == nil {
		logger.Infof("authentication failed")
		return nil
	}

	logger.Infof("user found by email")
	err = bcrypt.CompareHashAndPassword(usr.HashedPassword, []byte(password))
	if err != nil {
		logger.Errorf("authentication failed due to password", err)
		return nil
		// Todo: check what kind of error occurred
		//if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		//	return 0, ErrInvalidCredentials
		//} else {
		//	return 0, err
		//}
	}
	logger.Infof("authentication successful")
	return usr
}

func (s service) GenerateJWT(identity entity.UserAuthIdentity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       identity.GetID(),
		"email":    identity.GetEmail(),
		"username": identity.GetUserName(),
		"role":     identity.GetRole(),
		"exp":      time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
