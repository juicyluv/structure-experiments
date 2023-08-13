package domain

import (
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

func (p Post) Validate() error {
	if p.Title == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Title cannot be empty.", "title")
	}

	if p.Content == "" {
		return apperror.NewInvalidRequestError(ErrRequired, "Content cannot be empty.", "title")
	}

	return nil
}

type GetPostsFilters struct {
	AuthorID int64
	Page     int16
	Count    int16
}
