package postgres

import (
	"fmt"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, trigger bool, log *logger.Logger) error {

	if trigger {
		log.Info("📦 Migrating database...")
		models := []any{}
		schemas := []string{"user_module", "file_module"}

		log.Info("📦 Creating types...")

		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		for _, schema := range schemas {
			if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", schema)).Error; err != nil {
				return fmt.Errorf("failed to create schema %s: %w", schema, err)
			}
		}

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	log.Info("✅ Database connection successfully")
	return nil
}
