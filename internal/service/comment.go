package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/internal/service/apperror"
)

//go:generate mockgen -source=comment.go -destination=./mocks/comment.go -package=mocks
type CommentRepository interface {
	GetComments(ctx context.Context, postID int64) ([]domain.Comment, error)
	GetComment(ctx context.Context, id int64) (*domain.Comment, error)
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

		return nil, fmt.Errorf("getting post: %w", err)
	}

	comments, err := s.commentRepository.GetComments(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("getting comments for post %d: %w", postID, err)
	}

	return comments, nil
}

func (s *CommentService) AddComment(ctx context.Context, comment *domain.Comment) (int64, error) {
	err := comment.Validate()
	if err != nil {
		return 0, err
	}

	_, err = s.postRepository.GetPost(ctx, comment.PostID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, apperror.NewNotFoundError(err, "Post not found.", "postId")
		}

		return 0, fmt.Errorf("getting post: %w", err)
	}

	if comment.CommentID != nil {
		_, err = s.commentRepository.GetComment(ctx, *comment.CommentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return 0, apperror.NewNotFoundError(err, "Comment not found.", "commentId")
			}

			return 0, fmt.Errorf("getting comment: %w", err)
		}
	}

	id, err := s.commentRepository.AddComment(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("adding comment: %w", err)
	}

	return id, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	err := comment.Validate()
	if err != nil {
		return err
	}

	_, err = s.postRepository.GetPost(ctx, comment.PostID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NewNotFoundError(err, "Post not found.", "postId")
		}

		return fmt.Errorf("getting post: %w", err)
	}

	if comment.CommentID != nil {
		_, err = s.commentRepository.GetComment(ctx, *comment.CommentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.NewNotFoundError(err, "Comment not found.", "commentId")
			}

			return fmt.Errorf("getting comment: %w", err)
		}
	}

	err = s.commentRepository.UpdateComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("updating comment: %w", err)
	}

	return nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id int64) error {
	_, err := s.commentRepository.GetComment(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NewNotFoundError(err, "Comment not found.", "id")
		}

		return fmt.Errorf("getting comment: %w", err)
	}

	err = s.commentRepository.DeleteComment(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	return nil
}
