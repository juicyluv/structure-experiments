package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

type CommentRepository interface {
	GetComments(ctx context.Context, postID int64) ([]domain.Comment, error)
	AddComment(ctx context.Context, comment *domain.Comment) (int64, error)
	UpdateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
}

type CommentService struct {
	postRepository    PostRepository
	commentRepository CommentRepository
}

func NewCommentService(postRepo PostRepository, commentRepo CommentRepository) *CommentService {
	return &CommentService{postRepo, commentRepo}
}

func (s *CommentService) GetComments(ctx context.Context, postID int64) ([]domain.Comment, error) {
	_, err := s.postRepository.GetPost(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NewNotFoundError(err, "Post not found.", "id")
		}

		return nil, fmt.Errorf("getting post: %v", err)
	}

	comments, err := s.commentRepository.GetComments(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("getting comments for post %d: %v", postID, err)
	}

	return comments, nil
}
