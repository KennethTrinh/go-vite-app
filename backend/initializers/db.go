package initializers

import (
	"log"
	"os"
	"time"

	"github.com/KennethTrinh/go-vite-app/config"
	"github.com/KennethTrinh/go-vite-app/models"

	singlestore "github.com/singlestore-labs/gorm-singlestore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := config.Env.DatabaseUrl

	DB, err = gorm.Open(singlestore.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		// Logger: logger.Default.LogMode(logger.Silent),
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	// configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance! \n", err.Error())
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(10)           // Set the maximum number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(100)          // Set the maximum number of open connections to the database
	sqlDB.SetConnMaxLifetime(time.Hour) // Set the maximum amount of time a connection may be reused

	// DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	if !config.Env.Production {
		// DB = DB.Debug()
		log.Println("Running Migrations")
		err = DB.AutoMigrate(
			&models.Item{},
		)

		if err != nil {
			log.Fatal("Failed to run migrations! \n", err.Error())
			os.Exit(1)
		}
	}

	log.Println("🚀 Connected Successfully to the Database")

	stats := sqlDB.Stats()
	log.Printf("Open connections: %d, In-use connections: %d, Idle connections: %d\n",
		stats.OpenConnections, stats.InUse, stats.Idle)
}
