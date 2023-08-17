package postgresql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository/postgresql"
)

func TestCommentRepository_AddComment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := postgresql.NewCommentRepository(mock)

	q := "INSERT INTO comments"

	tests := []struct {
		name     string
		comment  *domain.Comment
		response int64
		wantErr  bool
		mock     func(comment *domain.Comment, id int64)
	}{
		{
			name: "args check",
			comment: &domain.Comment{
				AuthorID: 1,
				PostID:   2,
				Content:  "3",
			},
			response: 1,
			mock: func(comment *domain.Comment, id int64) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(id)

				mock.
					ExpectQuery(q).
					WithArgs(
						comment.AuthorID,
						comment.PostID,
						comment.Content,
						comment.CommentID,
					).WillReturnRows(rows)
			},
		},
		{
			name: "query error",
			comment: &domain.Comment{
				AuthorID: 1,
				PostID:   2,
				Content:  "3",
			},
			response: 0,
			mock: func(comment *domain.Comment, id int64) {
				mock.
					ExpectQuery(q).
					WithArgs(
						comment.AuthorID,
						comment.PostID,
						comment.Content,
						comment.CommentID,
					).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(tt *testing.T) {
			tc.mock(tc.comment, tc.response)

			resp, err := repo.AddComment(context.Background(), tc.comment)

			require.NoError(tt, mock.ExpectationsWereMet())

			require.Equal(tt, tc.response, resp)

			if !tc.wantErr {
				require.NoError(tt, err)
				return
			}

			require.Error(tt, err)
		})
	}
}
