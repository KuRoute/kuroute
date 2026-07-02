package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
)

type DB struct {
	*gorm.DB
}

func ConnectDB() (*DB, error) {
	env := godotenv.Load()
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// dsn, mean data source name. It is a string that contains the connection information for the database.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	// open a connection to postgresql
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// logger to take notes activiy query
		Logger: logger.Default.LogMode(logger.Info),
		// prepare statement to improve performance
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}


	log.Println("Connected to database")
	return &DB{DB: db}, nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}