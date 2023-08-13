package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

type (
	PostRepository interface {
		GetPost(ctx context.Context, id int64) (*domain.Post, error)
	}
)

type PostService struct {
	postRepository PostRepository
}

func NewPostService(postRepo PostRepository) *PostService {
	return &PostService{postRepo}
}

func (s *PostService) GetPost(ctx context.Context, id int64) (*domain.Post, error) {
	post, err := s.postRepository.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NewNotFoundError(err, "Post not found.", "id")
		}

		return nil, fmt.Errorf("getting post: %v", err)
	}

	return post, nil
}

func (s *PostService) AddPost(ctx context.Context, post *domain.Post) (int64, error) {
	return 0, nil
}
