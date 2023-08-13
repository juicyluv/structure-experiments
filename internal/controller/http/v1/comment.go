package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/juicyluv/structure-experiments/internal/controller/http/httpinout"
	"github.com/juicyluv/structure-experiments/internal/domain"
)

type CommentService interface {
	GetComments(ctx context.Context, postID int64) ([]domain.Comment, error)
}

type CommentController struct {
	commentService CommentService
}

func NewCommentController(s CommentService) CommentController {
	return CommentController{commentService: s}
}

func (c *CommentController) GetComments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpinout.BadRequest(err, "Invalid post ID.", "id", w, r)
		return
	}

	comments, err := c.commentService.GetComments(r.Context(), int64(id))
	if err != nil {
		httpinout.RespondWithError(err, w, r)
		return
	}

	httpinout.Json(w, comments, http.StatusOK)
}
