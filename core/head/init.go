package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/root9464/Go_GamlerDefi/config"
	"github.com/root9464/Go_GamlerDefi/database"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/root9464/Go_GamlerDefi/packages/middleware"
)

func (app *Core) init_http_server() {
	app.http_server = fiber.New()
	app.http_server.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))
	app.http_server.Use(middleware.LoggerMiddleware(app.logger))

	app.logger.Info("HTTP server initialized")
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

func (app *Core) init_routes() {
	api := app.http_server.Group("/api")
	app.modules.test.RegisterRoutes(api)
}
