package core

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/root9464/Go_GamlerDefi/config"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Core struct {
	config      *config.Config
	logger      *logger.Logger
	database    *mongo.Client
	http_server *fiber.App
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
		instance.init_http_server()
	})
	return instance
}
