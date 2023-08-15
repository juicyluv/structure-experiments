package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/juicyluv/structure-experiments/internal/domain"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository"
	"github.com/juicyluv/structure-experiments/internal/service"
	"github.com/juicyluv/structure-experiments/internal/service/mocks"
)

func comment(t *testing.T) (*service.CommentService, *mocks.MockCommentRepository, *mocks.MockPostRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	commentRepo := mocks.NewMockCommentRepository(mockCtl)
	postRepo := mocks.NewMockPostRepository(mockCtl)

	commentService := service.NewCommentService(postRepo, commentRepo)

	return commentService, commentRepo, postRepo
}

type test struct {
	name    string
	res     any
	err     error
	wantErr bool
}

func TestCommentService_GetComments(t *testing.T) {
	t.Parallel()

	srv, commRepo, postRepo := comment(t)

	tests := []struct {
		test
		postID int64
		mock   func(int64)
	}{
		{
			test: test{
				name: "empty result",
				res:  []domain.Comment{},
				err:  nil,
			},
			mock: func(postID int64) {
				postRepo.EXPECT().GetPost(context.Background(), postID).Return(&domain.Post{ID: 2}, nil)
				commRepo.EXPECT().GetComments(context.Background(), postID).Return([]domain.Comment{}, nil)
			},
			postID: 2,
		},
		{
			test: test{
				name: "successful response",
				res:  []domain.Comment{{ID: 1}, {ID: 2}},
				err:  nil,
			},
			mock: func(postID int64) {
				postRepo.EXPECT().GetPost(context.Background(), postID).Return(&domain.Post{ID: postID}, nil)
				commRepo.EXPECT().GetComments(context.Background(), postID).Return([]domain.Comment{{ID: 1}, {ID: 2}}, nil)
			},
			postID: 2,
		},
		{
			test: test{
				name:    "post not found",
				res:     []domain.Comment(nil),
				wantErr: true,
				err:     repository.ErrNotFound,
			},
			mock: func(postID int64) {
				postRepo.EXPECT().GetPost(context.Background(), postID).Return(nil, repository.ErrNotFound)
			},
			postID: 2,
		},
		{
			test: test{
				name:    "get post internal error",
				res:     []domain.Comment(nil),
				wantErr: true,
			},
			mock: func(postID int64) {
				postRepo.EXPECT().GetPost(context.Background(), postID).Return(nil, errors.New("some error"))
			},
			postID: 2,
		},
		{
			test: test{
				name:    "get comments internal error",
				res:     []domain.Comment(nil),
				wantErr: true,
			},
			mock: func(postID int64) {
				postRepo.EXPECT().GetPost(context.Background(), postID).Return(&domain.Post{ID: postID}, nil)
				commRepo.EXPECT().GetComments(context.Background(), postID).Return(nil, errors.New("some error"))
			},
			postID: 2,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tc.mock(tc.postID)

			res, err := srv.GetComments(context.Background(), tc.postID)

			require.Equal(t, res, tc.res)

			if !tc.wantErr {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}
