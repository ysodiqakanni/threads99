package post

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ysodiqakanni/threads99/internal/community"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	entity.Post
}

type CommentModel struct {
	entity.Comment
}

type Service interface {
	Get(ctx context.Context, id primitive.ObjectID) (Post, error)
	GetPostLiteById(ctx context.Context, id primitive.ObjectID) (*dto.PostResponse, error)
	CreatePost(ctx context.Context, request CreateNewPostRequest) (error, string)
	//AddCommentToPost(ctx context.Context, commentRequest AddCommentToPostRequest) error
	UpvoteComment(ctx context.Context, request CommentUpvoteRequest) error
	UpvotePost(ctx context.Context, request PostUpvoteRequest) error
	GetCommentsByPostId(ctx context.Context, postIdStr string) ([]entity.Comment, error)
	GetAllRecentPosts(ctx context.Context) ([]dto.TimelinePost, error)
	GetRecentPostsByCommunityId(ctx context.Context, communityId primitive.ObjectID) ([]dto.TimelinePost, error)
	GetRecentPostsByUserId(ctx context.Context, userId primitive.ObjectID) ([]dto.TimelinePost, error)
}

type service struct {
	repo          Repository
	communityRepo community.Repository
	logger        log.Logger
}

// NewService creates a new post service.
func NewService(repo Repository, communityRepo community.Repository, logger log.Logger) Service {
	return service{repo, communityRepo, logger}
}

type CreateNewPostRequest struct {
	Title             string `json:"title"`
	Content           string `json:"content"`
	CreatedByUserId   string `json:"created_by_user_id"`
	CreatedByUserName string
	CommunityId       string `json:"community_id"`
	CommunityName     string `json:"CommunityName"`
	PostType          string `json:"postType"`
}

type AddCommentToPostRequest struct {
	PostId          string `json:"post_id"`
	CreatedByUserId string `json:"created_by_user_id"`
	CommentContent  string `json:"comment_content"`
}

type PostUpvoteRequest struct {
	PostId    string `json:"post_id"`
	VoteValue int    `json:"vote_value"`
}

type CommentUpvoteRequest struct {
	PostId    string `json:"post_id"`
	CommentId string `json:"comment_id"`
	VoteValue int    `json:"vote_value"`
}

func (m CreateNewPostRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 256)),
		validation.Field(&m.Content, validation.Required, validation.Length(0, 1024)),
		validation.Field(&m.CommunityId, validation.Required),
	)
}
func (m AddCommentToPostRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CommentContent, validation.Required, validation.Length(0, 1024)),
		validation.Field(&m.PostId, validation.Required),
	)
}
func (m CommentUpvoteRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CommentId, validation.Required),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.VoteValue, validation.Required),
	)
}

func (m PostUpvoteRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.VoteValue, validation.Required),
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

func (s service) GetPostLiteById(ctx context.Context, id primitive.ObjectID) (*dto.PostResponse, error) {
	post, err := s.repo.GetPostLiteById(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s service) CreatePost(ctx context.Context, request CreateNewPostRequest) (error, string) {
	userId, err := primitive.ObjectIDFromHex(request.CreatedByUserId)

	if err != nil {
		return err, ""
	}
	communityId, err := primitive.ObjectIDFromHex(request.CommunityId)
	if err != nil {
		return err, ""
	}
	// now let's get community by ID
	community, err := s.communityRepo.Get(ctx, communityId)
	if err != nil {
		// error retrieving community object
		return err, ""
	}
	if community.Name == "" {
		return errors.New("The community with this ID cannot be found."), ""
	}

	post := entity.Post{
		Title:           request.Title,
		Content:         request.Content,
		CreatedByUserId: userId,
		CommunityID:     communityId,
		CommunityName:   request.CommunityName,
		Author: entity.Author{
			ID:       userId,
			Username: request.CreatedByUserName,
		},
		Metadata: entity.Metadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MainContent: entity.Content{
			Type: "Text", // Todo: change to enum and handle different post types.
			Body: request.Content,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, post)
	return err, id.Hex()
}

func (s service) GetAllRecentPosts(ctx context.Context) ([]dto.TimelinePost, error) {
	posts, err := s.repo.GetAllRecentPosts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (s service) GetRecentPostsByCommunityId(ctx context.Context, communityId primitive.ObjectID) ([]dto.TimelinePost, error) {
	posts, err := s.repo.GetRecentPostsByCommunityId(ctx, communityId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (s service) GetRecentPostsByUserId(ctx context.Context, userId primitive.ObjectID) ([]dto.TimelinePost, error) {
	posts, err := s.repo.GetRecentPostsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (s service) UpvotePost(ctx context.Context, request PostUpvoteRequest) error {
	postId, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return err
	}
	err = s.repo.UpvotePost(ctx, postId, request.VoteValue)
	return err
}

/*
func (s service) AddCommentToPost(ctx context.Context, commentRequest AddCommentToPostRequest) error {
	commentUserId, err := primitive.ObjectIDFromHex(commentRequest.CreatedByUserId)
	if err != nil {
		return err
	}
	comment := entity.Comment{
		ID:              primitive.NewObjectID(),
		Content:         commentRequest.CommentContent,
		CreatedByUserId: commentUserId,
		CreatedAt:       time.Now(),
	}
	postId, err := primitive.ObjectIDFromHex(commentRequest.PostId)
	if err != nil {
		return err
	}

	err = s.repo.AddCommentToPost(ctx, postId, comment)

	return err
}
*/

func (s service) UpvoteComment(ctx context.Context, request CommentUpvoteRequest) error {
	commentId, err := primitive.ObjectIDFromHex(request.CommentId)
	if err != nil {
		return err
	}
	postId, err := primitive.ObjectIDFromHex(request.PostId)
	if err != nil {
		return err
	}
	err = s.repo.UpvoteComment(ctx, commentId, postId, request.VoteValue)
	return err
}

func (s service) GetCommentsByPostId(ctx context.Context, postIdStr string) ([]entity.Comment, error) {
	postId, err := primitive.ObjectIDFromHex(postIdStr)
	emptyComments := []entity.Comment{}
	if err != nil {
		return emptyComments, err
	}
	results, err := s.repo.GetCommentsByPostId(ctx, postId)

	return results, err
}
