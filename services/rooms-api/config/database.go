package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitMySQL() *gorm.DB {
	_ = godotenv.Load()

	host := getenvOrDefault("DB_HOST", "localhost")
	port := getenvOrDefault("DB_PORT", "3306")
	user := getenvOrDefault("DB_USER", "root")
	password := getenvOrDefault("DB_PASSWORD", "root")
	dbName := getenvOrDefault("DB_NAME", "roomsdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error pinging MySQL: %v", err)
	}

	log.Printf("Connected to MySQL: %s@%s:%s/%s", user, host, port, dbName)

	return db
}
