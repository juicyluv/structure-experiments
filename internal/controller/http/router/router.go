package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	v1 "github.com/juicyluv/structure-experiments/internal/controller/http/v1"
)

func NewRouter(
	postController v1.PostController,
	commentController v1.CommentController,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/api/v1/posts/{id}", postController.GetPost)
	r.Get("/api/v1/posts", postController.GetPosts)
	r.Post("/api/v1/posts", postController.AddPost)
	r.Put("/api/v1/posts", postController.UpdatePost)

	r.Get("/api/v1/posts/{id}/comments", commentController.GetComments)

	return r
}
