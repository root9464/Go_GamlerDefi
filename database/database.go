package database

import (
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDatabase(url string, logger *logger.Logger) (*mongo.Client, error) {
	logger.Info("Connecting to MongoDB...")

	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		logger.Error("❌ Failed to connect to MongoDB")
		return nil, err
	}

	logger.Success("✅ Database connection successfully")
	return client, nil
}
