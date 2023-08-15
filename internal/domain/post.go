package domain

import (
	"fmt"
	"time"

	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

type Post struct {
	ID       int64
	AuthorID int64
	Title    string
	Content  string

	CreatedAt time.Time
	UpdatedAt *time.Time
}

const (
	postTitleMaxLength   = 128
	postContentMaxLength = 4096
)

func (p Post) Validate() error {
	if p.Title == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Title cannot be empty.", "title")
	}

	if len(p.Title) > postTitleMaxLength {
		return apperror.NewInvalidRequestError(
			ErrRequired,
			fmt.Sprintf("Title must be less than %d characters.", postTitleMaxLength),
			"title",
		)
	}

	if p.Content == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Content cannot be empty.", "content")
	}

	if len(p.Content) > postContentMaxLength {
		return apperror.NewInvalidRequestError(
			ErrRequired,
			fmt.Sprintf("Content must be less than %d characters.", postContentMaxLength),
			"content",
		)
	}

	return nil
}

type GetPostsFilters struct {
	AuthorID     int64
	Page         int
	PostsPerPage int
}
