package api

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/Ewan-Greer09/finance-app/api/config"
	"github.com/Ewan-Greer09/finance-app/api/database"
	"github.com/Ewan-Greer09/finance-app/api/handlers"
)

//go:embed web/*
var webFS embed.FS

type API struct {
	Name string

	*http.Server
	*slog.Logger
	config.Config
	*Handler
	ExpenseHandler *handlers.ExpenseHandler
	IncomeHandler  *handlers.IncomeHandler
	AdminHandler   *handlers.AdminHandler
}

func NewAPI() *API {
	cfg := config.LoadConfig()
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(cfg.API.LogLevel)}))

	api := &API{
		Name: cfg.API.ServiceName,
		Server: &http.Server{
			Addr: cfg.API.Addr,
		},
		Logger:         log,
		Config:         cfg,
		Handler:        NewHandler(log, cfg),
		ExpenseHandler: handlers.NewExpenseHandler(log, database.NewDatabase(cfg), webFS),
		IncomeHandler:  handlers.NewIncomeHandler(log, database.NewDatabase(cfg), webFS),
		AdminHandler:   handlers.NewAdminHandler(log, database.NewDatabase(cfg), webFS),
	}
	api.Server.Handler = api.registerRoutes()
	return api
}

func (a *API) Run() error {
	a.Info("Starting API server", "name", a.Name, "port", a.Server.Addr)

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			} else {
				a.Error("API server stopped with error", "error", err)
				return
			}
		}
	}()

	<-doneCh

	a.Info("Stopping API server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Error("API server shutdown failed", "error", err)
		return err
	}

	err := a.Handler.Database.Close()
	if err != nil {
		a.Error("Error while closing DB connection", "error", err)
		return err
	}

	a.Info("API server stopped")

	return nil
}

func (a *API) registerRoutes() http.HandlerFunc {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(a.Config.API.Timeout) * time.Second))

	// we only use templating on fragments, so we can serve the web directory directly
	// HTMX is used to make requests to the API, where we will use templates
	// Page navigation is handled with standard HTML, loading components with HTMX
	fs := http.FileServer(http.Dir("api/web"))
	r.Get("/*", http.StripPrefix("/", fs).ServeHTTP)

	r.Route("/admin", func(r chi.Router) {
		r.Use(a.AdminHandler.IsAdmin)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/web/admin.html")
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/expense", a.ExpenseHandler.Routes)
			r.Route("/income", a.IncomeHandler.Routes)
			r.Route("/admin", a.AdminHandler.Routes)
			r.Get("/graph", a.HandleGetExpensesAndIncomesGraph)
		})
	})

	return r.ServeHTTP
}
