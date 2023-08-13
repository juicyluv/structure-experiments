package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/pkg/postgres"
)

type PostRepository struct {
	*postgres.Postgres
}

func NewPostRepository(pg *postgres.Postgres) *PostRepository {
	return &PostRepository{pg}
}

func (r *PostRepository) GetPost(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post

	err := r.Pool.QueryRow(ctx, `
		SELECT 
		    id, author_id, title, content, created_at, updated_at
		FROM posts
		WHERE id=$1`, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("scanning port: %w", err)
	}

	return &post, nil
}
