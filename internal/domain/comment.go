package domain

import (
	"time"

	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

type Comment struct {
	ID        int64
	AuthorID  int64
	PostID    int64
	Content   string
	CommentID *int64

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewComment(authorID int64, postID int64, content string, commentID *int64) (*Comment, error) {
	if content == "" {
		return nil, apperror.NewInvalidRequestError(ErrRequired, "Content is required.", "content")
	}

	return &Comment{
		AuthorID:  authorID,
		PostID:    postID,
		Content:   content,
		CommentID: commentID,
	}, nil
}
