package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LullNil/authx-go/config"
	domainUser "github.com/LullNil/authx-go/domain/user"
	"github.com/LullNil/authx-go/internal/delivery/http/user"
	"github.com/LullNil/authx-go/internal/lib/logger"
	"github.com/LullNil/authx-go/internal/repository/postgres"
	users "github.com/LullNil/authx-go/internal/service/user"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/sync/errgroup"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type Services struct {
	User domainUser.Service
}

// Run starts the application.
func Run(cfg *config.Config) error {
	log := setupLogger(cfg.Env)

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Init database
	db, err := initPostgres(cfg, log)
	if err != nil {
		return err
	}
	defer func() {
		log.Debug("closing postgres connection...")
		db.Close()
	}()

	// Init app services
	appServices := initAppServices(db, log)

	// Init router
	router := initRouter(log, appServices)

	// Create errgroup for managing server goroutines
	group, gCtx := errgroup.WithContext(ctx)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:         cfg.HTTPServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	group.Go(func() error {
		log.Info("starting http server...", slog.String("port", cfg.HTTPServer.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server stopped with error", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	// Goroutine for handling shutdown signal and performing graceful shutdown
	group.Go(func() error {
		<-gCtx.Done()

		log.Debug("shutting down http server...")
		shutdownHTTPCtx, cancelHTTP := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelHTTP()
		if err := httpServer.Shutdown(shutdownHTTPCtx); err != nil {
			log.Error("http server graceful shutdown failed", slog.String("error", err.Error()))
		} else {
			log.Info("http server gracefully stopped")
		}

		return nil
	})

	// Wait for all goroutines to complete or for the first error
	if err := group.Wait(); err != nil {
		log.Error("application exited with error", slog.String("error", err.Error()))
		return err
	}

	log.Info("all servers gracefully stopped")
	return nil
}

// initPostgres connects to the database.
func initPostgres(cfg *config.Config, log *slog.Logger) (*sql.DB, error) {
	return postgres.ConnectWithRetries(
		context.Background(),
		cfg.Postgres,
		log,
	)
}

// initAppServices initializes the application services.
func initAppServices(db *sql.DB, log *slog.Logger) *Services {
	// Init repositories
	userRepo := postgres.NewUserRepository(db)

	// Init services
	userSvc := users.NewService(userRepo, log)

	return &Services{
		User: userSvc,
	}
}

// initRouter initializes the router.
func initRouter(log *slog.Logger, services *Services) http.Handler {
	// Init handlers
	userHandler := user.New(services.User, log)

	// Setup router
	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Recoverer)
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: add allowed origins
		AllowCredentials: true,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	}).Handler)

	// User routes
	router.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUser)
		r.Post("/login", userHandler.LoginUser)
		// r.Get("/info", userHandler.GetUserInfo)
	})

	return router
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(logger.NewPrettyHandler(os.Stdout, slog.LevelDebug))
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
}
