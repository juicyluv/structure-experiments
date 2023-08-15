package postgresql

import (
	"context"
	"fmt"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/pkg/postgres"
)

type CommentRepository struct {
	*postgres.Postgres
}

func NewCommentRepository(pg *postgres.Postgres) *CommentRepository {
	return &CommentRepository{pg}
}

func (r *CommentRepository) GetComments(ctx context.Context, postID int64) ([]domain.Comment, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT 
		    id, author_id, post_id, content, comment_id, created_at, updated_at
		FROM posts
		WHERE id=$1`, postID)
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
			return nil, fmt.Errorf("scanning port: %w", err)
		}
	}

	return comments, nil
}

func (r *CommentRepository) AddComment(ctx context.Context, comment *domain.Comment) (int64, error) {
	var id int64

	err := r.Pool.QueryRow(ctx, `
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
	_, err := r.Pool.Exec(ctx, `
		UPDATE comments
		SET content=$2, updated_at=now()
		WHERE id=$1`, comment.ID, comment.Content)
	if err != nil {
		return fmt.Errorf("updating comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	return nil
}
