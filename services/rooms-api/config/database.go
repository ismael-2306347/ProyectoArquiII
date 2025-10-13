package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitMongoDB() *mongo.Database {

	_ = godotenv.Load()

	uri := getenvOrDefault("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getenvOrDefault("MONGODB_DB", "roomsdb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ Error connecting to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Error pinging MongoDB: %v", err)
	}

	log.Printf("✅ Connected to MongoDB: %s", uri)
	log.Printf("✅ Using database: %s", dbName)

	return client.Database(dbName)
}
