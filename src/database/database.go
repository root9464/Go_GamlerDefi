package database

import (
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDatabase(url string, logger *logger.Logger, db_name string) (*mongo.Client, *mongo.Database, error) {
	logger.Info("Connecting to MongoDB...")

	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		logger.Error("❌ Failed to connect to MongoDB")
		return nil, nil, err
	}

	database := client.Database(db_name)

	logger.Success("✅ Database connection successfully")
	return client, database, nil
}
