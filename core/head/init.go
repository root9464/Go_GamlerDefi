package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/root9464/Go_GamlerDefi/config"
	"github.com/root9464/Go_GamlerDefi/database"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/root9464/Go_GamlerDefi/packages/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	gqlgen "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
)

func (app *Core) init_gql_server() {
	app.gql_server = fiber.New()
	app.gql_server.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))
	app.gql_server.Use(middleware.LoggerMiddleware(app.logger))

	srv := handler.New(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	app.gql_server.Get("/playground", adaptor.HTTPHandlerFunc(playground.Handler("Graphql Playground", "/query")))
	app.gql_server.All("/query", adaptor.HTTPHandler(srv))

	app.logger.Info("HTTP server initialized")
	app.logger.Successf("HTTP server listening on %s", app.config.Address())
	if err := app.gql_server.Listen(app.config.Address()); err != nil {
		app.logger.Errorf("Failed to start HTTP server: %v", err)
	}
}

func (app *Core) init_database() {
	if app.config == nil {
		app.logger.Error("Config is not initialized, cannot connect to database")
		return
	}
	mdb, err := database.ConnectDatabase(app.config.DatabaseUrl, app.logger)
	if err != nil {
		app.logger.Errorf("Failed to connect to database: %v", err)
	}

	app.database = mdb

}

func (app *Core) init_logger() {
	if app.logger == nil {
		app.logger = logger.GetLogger()
	}
}

func (app *Core) init_config() {
	if app.config == nil {
		config, err := config.LoadConfig("../.env")
		if err != nil {
			app.logger.Errorf("Failed to load config: %v", err)
		}

		app.config = config
	}
}
