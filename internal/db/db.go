package db

import (
	"log"
	"os"
	"url_shortener/configs"
	"url_shortener/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

func ConnectDB(config *configs.Config, logger logger.Logger) *gorm.DB {
	gormConfig := gorm.Config{}
	if config.ENVIRONMENT != configs.DEV {
		gormConfig.Logger = gorm_logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gorm_logger.Config{
				LogLevel:             gorm_logger.Silent,
				ParameterizedQueries: true,
				Colorful:             false,
			},
		)
	}
	db, err := gorm.Open(postgres.Open(config.DB.POSTGRES_DSN), &gormConfig)

	logger.Debug("Successfully connect to db: ", db.Name())
	logger.Info("Successfully connect to db: ", db.Name())
	if err != nil {
		logger.Error(err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Shortener{})
	return db
}
