package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ysodiqakanni/threads99/internal/entity"
)

type CreateNewCommentRequest struct {
	PostId          string
	ParentId        string
	ContentText     string
	MediaUrls       []string
	CreatedByUserId string
}

// Either of the comment text and media is required.
func (m CreateNewCommentRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreatedByUserId, validation.Required),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(
			&m.ContentText,
			validation.When(len(m.MediaUrls) == 0, validation.Required).Else(validation.Empty),
		),
		validation.Field(
			&m.MediaUrls,
			validation.When(m.ContentText == "", validation.Required).Else(validation.Empty),
		),
	)
}

// A tree of nested comments.
type CommentTree struct {
	Comment entity.Comment `json:"comment"`
	Replies []CommentTree  `json:"replies"`
}
