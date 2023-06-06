package initializers

import (

	"log"
	"os"

	"github.com/wpcodevo/golang-fiber/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {

	connStr := "user=shahzodsharifov password=qnNurw9t8QYK dbname=neondb host=ep-winter-mouse-885421.eu-central-1.aws.neon.tech sslmode=verify-full"
	var err error
	// dsn :=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err !=nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.UserRelationship{})
	DB.AutoMigrate(&models.Post{})
	DB.AutoMigrate(&models.Comment{})
	DB.AutoMigrate(&models.Vaccancy{})
	DB.AutoMigrate(&models.Event{})
	DB.AutoMigrate(&models.Like{})

	log.Println("Connected successfully to the database")
}