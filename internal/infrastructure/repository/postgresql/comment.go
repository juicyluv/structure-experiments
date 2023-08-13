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
		return nil, fmt.Errorf("selecting comments: %v", err)
	}
	defer rows.Close()

	var (
		comment  domain.Comment
		comments []domain.Comment
	)
	for rows.Next() {
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
