package rest

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"github.com/rs/cors"

	// Need this for all repositories
	_ "github.com/go-sql-driver/mysql"

	betshandler "github.com/MartinNikolovMarinov/sb-bets/handlers/bets-handler"
	betsrepo "github.com/MartinNikolovMarinov/sb-bets/repositories/bets-repo"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-infra/pubsub"
)

// App comment
type App struct {
	Router    *chi.Mux
	BetsRepo  betsrepo.BetsRepo
	Validator *validator.Validate
}

// Init method initializes the App
func (a *App) Init() {
	var connectionString = "user:password@tcp(127.0.0.1:3306)/sports_betting_db"
	args := os.Args[1:]
	if len(args) > 0 {
		connectionString = args[0]
	}
	a.BetsRepo = betsrepo.NewMysqlRepo(connectionString)
	a.Validator = validator.New()
	a.Router = chi.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(pubsub.MetricsMiddleware)

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
			Get("/bets/{userID}", betshandler.All(a.BetsRepo, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Get("/bets/all", betshandler.All(a.BetsRepo, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Post("/bets/place", betshandler.Place(a.BetsRepo, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Post("/bets/payout", betshandler.Payout(a.BetsRepo, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Delete("/bets/cancel", betshandler.Cancel(a.BetsRepo, a.Validator).ServeHTTP)
	})
}

// Run starts the REST API server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
