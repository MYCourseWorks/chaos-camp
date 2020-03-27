package rest

import (
	"log"
	"net/http"
	"os"

	gameshandler "github.com/MartinNikolovMarinov/sb-feed/handlers/games-handler"
	leaguehandler "github.com/MartinNikolovMarinov/sb-feed/handlers/leagues-handler"
	sportshandler "github.com/MartinNikolovMarinov/sb-feed/handlers/sports-handler"
	gamesrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/games-repo"
	leaguesrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/leagues-repo"
	sportsrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/sports-repo"
	"github.com/MartinNikolovMarinov/sb-feed/scrape"
	oddsapisrc "github.com/MartinNikolovMarinov/sb-feed/scrape/sources/odds-api"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-infra/pubsub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator"
	"github.com/rs/cors"

	// Need this for all repositories
	_ "github.com/go-sql-driver/mysql"
)

// App comment
type App struct {
	Router      *chi.Mux
	SportsRepo  sportsrepo.SportsRepo
	LeaguesRepo leaguesrepo.LeaguesRepo
	GamesRepo   gamesrepo.GamesRepo
	Validator   *validator.Validate
	feedWorker  *scrape.FeedWorker
}

// Init method initializes the App
func (a *App) Init() {
	var connectionString = "user:password@tcp(127.0.0.1:3306)/sports_betting_db"
	args := os.Args[1:]
	if len(args) > 0 {
		connectionString = args[0]
	}
	a.Validator = validator.New()
	a.SportsRepo = sportsrepo.NewMysqlRepo(connectionString)
	a.LeaguesRepo = leaguesrepo.NewMysqlRepo(connectionString)
	a.GamesRepo = gamesrepo.NewMysqlRepo(connectionString)
	a.Router = chi.NewRouter()
	a.feedWorker = scrape.NewFeedWorker(connectionString)
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
			Get("/sports/all", sportshandler.All(a.SportsRepo).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Get("/leagues/all", leaguehandler.All(a.LeaguesRepo).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Get("/games/all", gameshandler.All(a.GamesRepo, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Get("/games/freeze/{eventID}", gameshandler.Freeze(a.GamesRepo, a.Validator).ServeHTTP)
	})
}

// Run starts the REST API server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Scrape comment
func (a *App) Scrape() {
	httpClient := infra.NewFSClient()
	feed1 := oddsapisrc.New(httpClient, "./scrape/sources/odds-api/raw-data/")
	a.feedWorker.AddFeed(feed1)
	a.feedWorker.Work()
}
