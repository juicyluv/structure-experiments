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
			ErrInvalidValue,
			fmt.Sprintf("Title must be less than %d characters.", postTitleMaxLength),
			"title",
		)
	}

	if p.Content == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Content cannot be empty.", "content")
	}

	if len(p.Content) > postContentMaxLength {
		return apperror.NewInvalidRequestError(
			ErrInvalidValue,
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

func (f *GetPostsFilters) Validate() error {
	if f.Page < 1 {
		return apperror.NewInvalidRequestError(
			ErrInvalidValue,
			"Page must be greater than zero.",
			"page",
		)
	}

	if f.PostsPerPage < 1 {
		return apperror.NewInvalidRequestError(
			ErrInvalidValue,
			"Posts per page must be greater than zero.",
			"postsPerPage",
		)
	}

	return nil
}
