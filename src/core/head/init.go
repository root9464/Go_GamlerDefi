package core

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/root9464/Go_GamlerDefi/src/config"
	"github.com/root9464/Go_GamlerDefi/src/database"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	slog_logger "github.com/root9464/Go_GamlerDefi/src/packages/lib/slog_logger"
	"github.com/root9464/Go_GamlerDefi/src/packages/middleware"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	_ "github.com/root9464/Go_GamlerDefi/docs"
)

func (app *Core) init_http_server() {
	app.http_server = fiber.New()

	// app.http_server.Use(cors.New(cors.Config{
	// 	AllowOrigins: strings.Join([]string{
	// 		"https://gamler.atma-dev.ru",
	// 		"https://serv.gamler.online",
	// 		"https://gamler.online",
	// 	}, ","),
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization, Wallet-Address",
	// 	AllowMethods: strings.Join([]string{
	// 		fiber.MethodGet,
	// 		fiber.MethodPost,
	// 		fiber.MethodHead,
	// 		fiber.MethodPut,
	// 		fiber.MethodDelete,
	// 		fiber.MethodPatch,
	// 	}, ","),
	// 	AllowCredentials: false,
	// }))

	app.http_server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	app.http_server.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Next()
	})

	app.http_server.Use(middleware.LoggerMiddleware(app.logger))
	app.http_server.Use(middleware.ErrorMiddleware)

	app.logger.Info("HTTP server initialized")
}

func (app *Core) init_database() {
	if app.config == nil {
		app.logger.Error("Config is not initialized, cannot connect to database")
		return
	}
	client, database, err := database.ConnectDatabase(app.config.DatabaseUrl, app.logger, app.config.DatabaseName)
	if err != nil {
		app.logger.Errorf("Failed to connect to database: %v", err)
	}

	app.db_client = client
	app.database = database
}

func (app *Core) init_logger() {
	if app.logger == nil {
		app.logger = logger.GetLogger()
	}
}

func (app *Core) init_slog_logger() {
	if app.slog_logger == nil {
		app.slog_logger = slog_logger.NewLogger(&slog_logger.Option{
			Colorize:         app.config.Logger.Colorize,
			Level:            app.config.Logger.Level,
			Format:           app.config.Logger.Format,
			OutputEmptyAttrs: app.config.Logger.OutputEmptyAttrs,
			TimeFormat:       app.config.Logger.TimeFormat,

			Source: slog_logger.Source{
				Add:        app.config.Logger.Source.Add,
				ShowLine:   app.config.Logger.Source.ShowLine,
				TrimPrefix: app.config.Logger.Source.TrimPrefix,
				ShowFunc:   app.config.Logger.Source.ShowFunc,
				PathMode:   app.config.Logger.Source.PathMode,
			},
		})
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

func (app *Core) init_validator() {
	if app.validator == nil {
		app.validator = validator.New()
	}
}

func (app *Core) init_ton_client() {
	client := liteclient.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(context.Background(), app.config.TonConnect)
	if err != nil {
		app.logger.Errorf("Failed to add connections from config url: %v", err)
	}
	app.ton_client = ton.NewAPIClient(client)

	app.logger.Info("💎 TON client initialize successfully")
}

func (app *Core) init_ton_api() {
	client, err := tonapi.NewClient(tonapi.TonApiURL, &tonapi.Security{})
	if err != nil {
		app.logger.Errorf("Failed to create ton api client: %v", err)
	}
	app.ton_api = client

	app.logger.Info("🔷 TON api initialize successfully")
}

func (app *Core) init_routes() {
	app.http_server.Get("/web3/swagger/*", swagger.HandlerDefault)
	api := app.http_server.Group("/api")
	app.modules.test.RegisterRoutes(api)
	app.modules.referral.RegisterRoutes(api)
	app.modules.validation.RegisterRoutes(api)
	app.modules.ton.RegisterRoutes(api)
	app.modules.game_hub.InitDelivery(api)
}
