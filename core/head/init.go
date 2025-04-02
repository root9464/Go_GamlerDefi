package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/root9464/Go_GamlerDefi/config"
	"github.com/root9464/Go_GamlerDefi/database"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/root9464/Go_GamlerDefi/modules/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

func (app *Core) init_http_server() {
	app.http_server = fiber.New()
	app.http_server.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	app.http_server.Get("/playground", adaptor.HTTPHandlerFunc(playground.Handler("Graphql Playground", "/query")))
	app.http_server.All("/query", adaptor.HTTPHandler(srv))

	app.logger.Info("HTTP server initialized")
	app.logger.Successf("HTTP server listening on %s", app.config.Address())
	app.http_server.Listen(app.config.Address())
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

func (app *Core) init_logger() error {
	if app.logger == nil {
		app.logger = logger.GetLogger()
	}
	return nil
}

func (app *Core) init_config() {
	if app.config == nil {
		config, err := config.LoadConfig("../.env")
		if err != nil {

		}

		app.config = config
	}
}
