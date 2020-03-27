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

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-infra/pubsub"
	sitehandler "github.com/MartinNikolovMarinov/sb-users/handlers/site-handler"
	userhandler "github.com/MartinNikolovMarinov/sb-users/handlers/user-handler"
	siterepo "github.com/MartinNikolovMarinov/sb-users/repositories/site-repo"
	userrepo "github.com/MartinNikolovMarinov/sb-users/repositories/user-repo"
)

// App comment
type App struct {
	Router    *chi.Mux
	Users     userrepo.UserRepo
	Sites     siterepo.SiteRepo
	Validator *validator.Validate
}

// Init method initializes the App
func (a *App) Init() {
	var connectionString = "user:password@tcp(127.0.0.1:3306)/sports_betting_db"
	args := os.Args[1:]
	if len(args) > 0 {
		connectionString = args[0]
	}
	a.Users = userrepo.NewMysqlUserRepo(connectionString)
	a.Sites = siterepo.NewMysqlSiteRepo(connectionString)
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
		r.Post("/users/login", userhandler.Login(a.Users, a.Validator).ServeHTTP)
		r.Post("/users/register", userhandler.Register(a.Users, a.Sites, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Get("/users/all/{siteID}", userhandler.All(a.Users, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Delete("/users/delete/{userID}", userhandler.Delete(a.Users, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Post("/users/create", userhandler.Create(a.Users, a.Sites, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Post("/sites/create", sitehandler.Create(a.Sites, a.Validator).ServeHTTP)
		r.With(infra.JwtVerifyMiddleware).
			Put("/users/update", userhandler.Update(a.Users, a.Validator).ServeHTTP)
	})
}

// Run starts the REST API server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
