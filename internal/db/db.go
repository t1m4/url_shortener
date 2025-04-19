package db

import (
	"log"
	"os"
	"url_shortener/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(config *configs.Config) *gorm.DB {
	gormConfig := gorm.Config{}
	if config.ENVIRONMENT != configs.DEV {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:             logger.Silent,
				ParameterizedQueries: true,
				Colorful:             false,
			},
		)
	}
	db, err := gorm.Open(postgres.Open(config.POSTGRES_DSN), &gormConfig)

	log.Println("Successfully connect to db: ", db)
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Shortener{})
	return db
}
