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

func (r *PostRepository) AddPost(ctx context.Context, post *domain.Post) (int64, error) {
	var id int64

	err := r.Pool.QueryRow(ctx, `
		INSERT INTO posts(author_id, title, content)
		VALUES($1, $2, $3)
		RETURNING id`,
		post.AuthorID,
		post.Title,
		post.Content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting post: %w", err)
	}

	return id, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, filters *domain.GetPostsFilters) ([]domain.Post, error) {
	var (
		filter string
		argID  = 1
		args   []any
	)

	if filters.AuthorID != 0 {
		filter += "author_id = $1 "
		args = append(args, filters.AuthorID)
		argID++
	}

	if filter != "" {
		filter = "WHERE " + filter
	}

	rows, err := r.Pool.Query(ctx, fmt.Sprintf(`
		SELECT 
		    id, author_id, title, content, created_at, updated_at
		FROM posts
		%s
		%s`, filter, paginationQuery(filters.Page, filters.PostsPerPage)), args...)
	if err != nil {
		return nil, fmt.Errorf("selecting posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post

		err = rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning post: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE posts
		SET title=$2, content=$3, updated_at=now()
		WHERE id=$1`,
		post.ID,
		post.Title,
		post.Content,
	)
	if err != nil {
		return fmt.Errorf("updating post: %w", err)
	}

	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id int64) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting post: %w", err)
	}

	return nil
}
