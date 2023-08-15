package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juicyluv/structure-experiments/internal/config"
	"github.com/juicyluv/structure-experiments/internal/controller/http/router"
	v1 "github.com/juicyluv/structure-experiments/internal/controller/http/v1"
	"github.com/juicyluv/structure-experiments/internal/infrastructure/repository/postgresql"
	"github.com/juicyluv/structure-experiments/internal/pkg/logger"
	"github.com/juicyluv/structure-experiments/internal/service"
	"github.com/juicyluv/structure-experiments/pkg/httpserver"
	"github.com/juicyluv/structure-experiments/pkg/postgres"
)

type App struct {
}

func New() (*App, error) {

	return nil, nil
}

func Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cfg := config.Get()

	// Connect to Postgres
	pg, err := postgres.New(ctx, cfg.Repository.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("failed to create postgres: %v", err)
	}
	defer pg.Close()

	// Initialize repositories
	postRepo := postgresql.NewPostRepository(pg.Pool)
	commentRepo := postgresql.NewCommentRepository(pg.Pool)

	// Initialize services
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(postRepo, commentRepo)

	// Initialize controllers
	postController := v1.NewPostController(postService)
	commentController := v1.NewCommentController(commentService)

	r := router.NewRouter(postController, commentController)

	server := httpserver.New(
		r,
		httpserver.Port(cfg.HttpServer.Port),
		httpserver.ShutdownTimeout(time.Duration(cfg.HttpServer.ShutdownTimeout)*time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.HttpServer.WriteTimeout)*time.Second),
		httpserver.ReadTimeout(time.Duration(cfg.HttpServer.ReadTimeout)*time.Second),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Get().Info("Interrupt signal: " + s.String())
	case err = <-server.Notify():
		logger.Get().Error(fmt.Errorf("server stop signal received: %w", err))
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		logger.Get().Error(fmt.Errorf("failed to shutdown the server: %w", err))
	}
	logger.Get().Info("Server has been shut down successfully")

	return nil
}
