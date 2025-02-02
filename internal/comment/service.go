package comment

import (
	"context"
	"errors"
	"github.com/ysodiqakanni/threads99/internal/dto"
	"github.com/ysodiqakanni/threads99/internal/entity"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Service interface {
	CreateNewComment(ctx context.Context, request dto.CreateNewCommentRequest) (error, string)
}
type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new post service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) CreateNewComment(ctx context.Context, request dto.CreateNewCommentRequest) (error, string) {
	// Todo: required fields: postId, AuthorId and name, ContentBody or Media

	// PostId and userId are required.
	postId, err := primitive.ObjectIDFromHex(request.PostId)
	userId, err1 := primitive.ObjectIDFromHex(request.CreatedByUserId)
	if err != nil || err1 != nil {
		return errors.New("UserId and Post Id are required"), ""
	}
	var parentId primitive.ObjectID
	if request.ParentId != "" && &request.ParentId != nil {
		// parentId could be null
		parentId, err = primitive.ObjectIDFromHex(request.ParentId)
		if err != nil {
			return errors.New("Invalid Parent postId"), ""
		}
	}
	username := "testUsername" // Todo: retrieve from the jwt

	// Todo: Question: Should we check if the post and parentComment? exist?
	comment := entity.Comment{
		PostID:   postId,
		ParentID: &parentId,
		Content: entity.Content{
			Body:  request.ContentText,
			Media: request.MediaUrls,
		},
		Author: entity.Author{
			ID:       userId,
			Username: username,
		},
		Metadata: entity.CommentMetadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	err, id := s.repo.CreateNewComment(ctx, comment)
	return err, id.Hex()
}
