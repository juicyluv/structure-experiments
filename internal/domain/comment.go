package domain

import (
	"fmt"
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

const commentContentMaxLength = 512

func (c *Comment) Validate() error {
	if c.Content == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Content is required.", "content")
	}

	if len(c.Content) > commentContentMaxLength {
		return apperror.NewInvalidRequestError(
			ErrInvalidValue,
			fmt.Sprintf("Content must be less than %d characters.", commentContentMaxLength),
			"content",
		)
	}

	return nil
}
