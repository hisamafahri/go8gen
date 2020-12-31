package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"abc/configs"
	"abc/internal/domain/health"
	healthHTTP "abc/internal/domain/health/handler/http"
	// inject:import
bookHTTP "abc/internal/domain/book/handler/http"
bookPostgres "abc/internal/domain/book/repository/postgres"
bookUseCase "abc/internal/domain/book/usecase"
	"abc/internal/domain/health/repository/postgres"
	"abc/internal/domain/health/usecase"
	"abc/internal/middleware"
	"abc/third_party/database"
)

type App struct {
	httpServer *http.Server
	healthUC health.UseCase
	//inject:app
bookUC *bookUseCase.BookUseCase
}

func NewApp(cfg *configs.Configs) *App {
	db := database.NewSqlx(cfg)

	return &App{
		healthUC: usecase.NewHealthUseCase(postgres.NewHealthRepository(db)),
		// inject:usecase
bookUC: bookUseCase.NewBookUseCase(bookPostgres.NewBookRepository(db)),
	}
}

func (a *App) Run(cfg *configs.Configs, version string) error {
	router := chi.NewRouter()
	router.Use(middleware.Cors)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	healthHTTP.RegisterHTTPEndPoints(router, a.healthUC)
	// inject:handler
bookHTTP.RegisterHTTPEndPoints(router, a.bookUC)

	a.httpServer = &http.Server{
		Addr:           ":" + cfg.Api.Port,
		Handler:        router,
		ReadTimeout:    cfg.Api.ReadTimeout * time.Second,
		WriteTimeout:   cfg.Api.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("API version: %s\n", version)
		log.Printf("serving at %s:%s\n", cfg.Api.Host, cfg.Api.Port)
		printAllRegisteredRoutes(router)
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), cfg.Api.IdleTimeout*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func printAllRegisteredRoutes(router *chi.Mux) {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("path: %s method: %s ", path, method)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Print(err)
	}
}
