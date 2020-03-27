package rest

import (
	"log"
	"net/http"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-infra/pubsub"
	metricshandler "github.com/MartinNikolovMarinov/sb-metrics/handlers/metrics-handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"github.com/rs/cors"
)

// App comment
type App struct {
	Router       *chi.Mux
	MetricsCache pubsub.Cache
	Validator    *validator.Validate
}

// Init method initializes the App
func (a *App) Init() {
	c, err := pubsub.NewKafkaCache(nil, "Metrics", 0)
	if err != nil {
		panic("Could not load Metrics cache")
	}

	a.Validator = validator.New()
	a.MetricsCache = c
	a.Router = chi.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	var cors = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"*"},
	})
	a.Router.Use(cors.Handler)

	a.Router.Route("/api/v1", func(r chi.Router) {
		r.With(infra.JwtVerifyMiddleware).
			Get("/metrics/all", metricshandler.All(a.MetricsCache, a.Validator).ServeHTTP)
	})
}

// Run starts the REST API server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
