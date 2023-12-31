package postgresql

import (
	"context"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
)

type CommentRepository struct {
	driver Driver
}

func NewCommentRepository(dr Driver) *CommentRepository {
	return &CommentRepository{driver: dr}
}

func (r *CommentRepository) GetComments(ctx context.Context, postID int64) ([]domain.Comment, error) {
	rows, err := r.driver.Query(ctx, `
		SELECT 
		    id, author_id, post_id, content, comment_id, created_at, updated_at
		FROM comments
		WHERE post_id=$1`, postID)
	if err != nil {
		return nil, fmt.Errorf("selecting comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err = rows.Scan(
			&comment.ID,
			&comment.AuthorID,
			&comment.PostID,
			&comment.Content,
			&comment.CommentID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning comment: %w", err)
		}
	}

	return comments, nil
}

func (r *CommentRepository) GetComment(ctx context.Context, id int64) (*domain.Comment, error) {
	var comment domain.Comment

	err := r.driver.QueryRow(ctx, `
		SELECT 
		    id, author_id, post_id, content, comment_id, created_at, updated_at
		FROM comments
		WHERE id=$1`, id).Scan(
		&comment.ID,
		&comment.AuthorID,
		&comment.PostID,
		&comment.Content,
		&comment.CommentID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning comment: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepository) AddComment(ctx context.Context, comment *domain.Comment) (int64, error) {
	var id int64

	err := r.driver.QueryRow(ctx, `
		INSERT INTO comments(author_id, post_id, content, comment_id)
		VALUES($1,$2,$3,$4)
		RETURNING id`,
		comment.AuthorID,
		comment.PostID,
		comment.Content,
		comment.CommentID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting comment: %w", err)
	}

	return id, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	_, err := r.driver.Exec(ctx, `
		UPDATE comments
		SET content=$2, updated_at=now()
		WHERE id=$1`, comment.ID, comment.Content)
	if err != nil {
		return fmt.Errorf("updating comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int64) error {
	_, err := r.driver.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	return nil
}
