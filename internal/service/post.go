package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

//go:generate mockgen -source=post.go -destination=./mocks/post.go -package=mocks
type PostRepository interface {
	GetPost(ctx context.Context, id int64) (*domain.Post, error)
	GetPosts(ctx context.Context, filters *domain.GetPostsFilters) ([]domain.Post, error)
	AddPost(ctx context.Context, post *domain.Post) (int64, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, id int64) error
}

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

func (s *PostService) GetPosts(ctx context.Context, filters *domain.GetPostsFilters) ([]domain.Post, error) {
	err := filters.Validate()
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepository.GetPosts(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("getting post: %v", err)
	}

	return posts, nil
}

func (s *PostService) AddPost(ctx context.Context, post *domain.Post) (int64, error) {
	err := post.Validate()
	if err != nil {
		return 0, err
	}

	id, err := s.postRepository.AddPost(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("adding post: %w", err)
	}

	return id, nil
}

func (s *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	err := post.Validate()
	if err != nil {
		return err
	}

	_, err = s.postRepository.GetPost(ctx, post.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NewNotFoundError(err, "Post not found.", "id")
		}

		return fmt.Errorf("getting post %d: %w", post.ID, err)
	}

	err = s.postRepository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("updating post: %w", err)
	}

	return nil
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	post, err := s.postRepository.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NewNotFoundError(err, "Post not found.", "id")
		}

		return fmt.Errorf("getting post %d: %w", post.ID, err)
	}

	err = s.postRepository.DeletePost(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("deleting post: %w", err)
	}

	return nil
}
