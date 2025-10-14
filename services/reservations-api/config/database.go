package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"reservations-api/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitDatabase() *gorm.DB {
	// Carga .env si existe (útil en local; en Docker no hace falta)
	_ = godotenv.Load()

	user := getenvOrDefault("DB_USER", "user")
	pass := getenvOrDefault("DB_PASSWORD", "userpass")
	host := getenvOrDefault("DB_HOST", "mysql") // nombre del servicio en Compose
	port := getenvOrDefault("DB_PORT", "3306")
	name := getenvOrDefault("DB_NAME", "reservationsdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, port, name,
	)

	var (
		db  *gorm.DB
		err error
	)

	// Backoff exponencial simple (hasta ~1 min)
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, e := db.DB()
			if e != nil {
				err = e
			} else {
				if pingErr := sqlDB.Ping(); pingErr != nil {
					err = pingErr
				} else {
					log.Printf("Conectado a MySQL en %s:%s (intento %d)", host, port, attempt)
					goto MIGRATE
				}
			}
		}

		wait := time.Duration(attempt*2) * time.Second
		log.Printf("MySQL no listo (intento %d): %v. Reintentando en %s...", attempt, err, wait)
		time.Sleep(wait)
	}

	log.Fatalf("No se pudo conectar a la base tras reintentos: %v", err)

MIGRATE:
	if err := db.AutoMigrate(&domain.Reservation{}); err != nil {
		log.Fatalf("Error en AutoMigrate: %v", err)
	}
	log.Println("Migración de base de datos completada")

	return db
}
