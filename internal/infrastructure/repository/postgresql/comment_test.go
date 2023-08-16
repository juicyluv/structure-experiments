package postgresql_test

import (
	"context"
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

	tests := []struct {
		name     string
		comment  *domain.Comment
		response int64
		wantErr  bool
		err      error
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
					ExpectQuery("INSERT INTO comments").
					WithArgs(
						comment.AuthorID,
						comment.PostID,
						comment.Content,
						comment.CommentID,
					).WillReturnRows(rows)
			},
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
			require.ErrorIs(tt, err, tc.err)
		})
	}
}
