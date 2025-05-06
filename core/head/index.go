package core

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/config"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"go.mongodb.org/mongo-driver/mongo"
)

var client = liteclient.NewConnectionPool()

type Core struct {
	config      *config.Config
	logger      *logger.Logger
	database    *mongo.Client
	validator   *validator.Validate
	ton_client  *ton.APIClient
	ton_api     *tonapi.Client
	http_server *fiber.App
	modules     *Modules
}

var (
	instance *Core
	once     sync.Once
)

func InitApp() *Core {
	instance = &Core{}
	once.Do(func() {
		instance.init_config()
		instance.init_logger()
		instance.init_database()
		instance.init_validator()
		instance.init_ton_client()

		instance.init_http_server()
		instance.init_modules()
		instance.init_routes()
	})
	return instance
}

func (app *Core) Start() {
	app.logger.Successf("HTTP server listening on %s", app.config.Address())
	if err := app.http_server.Listen(app.config.Address()); err != nil {
		app.logger.Errorf("Failed to start HTTP server: %v", err)
	}
}
