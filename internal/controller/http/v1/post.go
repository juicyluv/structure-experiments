package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/juicyluv/structure-experiments/internal/controller/http/httpinout"
	"github.com/juicyluv/structure-experiments/internal/domain"
)

type PostService interface {
	GetPost(ctx context.Context, id int64) (*domain.Post, error)
	AddPost(ctx context.Context, post *domain.Post) (int64, error)
	//GetPosts(ctx context.Context, filters *domain.GetPostsFilters) ([]domain.Post, error)
	//UpdatePost(ctx context.Context, post *domain.Post) error
}

type PostController struct {
	postService PostService
}

func NewPostController(s PostService) PostController {
	return PostController{postService: s}
}

func (c *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpinout.BadRequest(err, "Invalid post ID.", "", w, r)
		return
	}

	post, err := c.postService.GetPost(r.Context(), int64(id))
	if err != nil {
		httpinout.RespondWithError(err, w, r)
		return
	}

	httpinout.Json(w, post, http.StatusOK)
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {

}

func (c *PostController) AddPost(w http.ResponseWriter, r *http.Request) {
	type request struct {
		AuthorID int64
		Title    string
		Content  string
	}

	var req request
	err := httpinout.ReadJSON(r, &req)
	if err != nil {
		httpinout.JsonError(err, w, r)
		return
	}

	post := &domain.Post{
		AuthorID: req.AuthorID,
		Title:    req.Title,
		Content:  req.Content,
	}

	id, err := c.postService.AddPost(r.Context(), post)
	if err != nil {
		httpinout.RespondWithError(err, w, r)
		return
	}

	type response struct {
		PostID int64 `json:"post_id"`
	}

	httpinout.Json(w, response{PostID: id}, http.StatusOK)
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {

}
