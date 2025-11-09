package main

import (
	"log"
	"os"
	"search-api/config"
	"search-api/controllers"
	"search-api/events"
	"search-api/repositories"
	"search-api/services"
	"time"

	"github.com/gin-gonic/gin"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Configuración
	solrURL := getEnv("SOLR_URL", "http://localhost:8983/solr/rooms_core")
	memcachedHost := getEnv("MEMCACHED_HOST", "localhost")
	memcachedPort := getEnv("MEMCACHED_PORT", "11211")
	rabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	_ = rabbitMQURL

	// Inicializar Solr
	solrClient := config.NewSolrClient(solrURL)

	// Inicializar Memcached
	cacheRepo := repositories.NewCacheRepository(
		memcachedHost,
		memcachedPort,
		5*time.Minute,
	)

	// Inicializar RabbitMQ
	rabbitConn := config.InitRabbitMQ()
	if rabbitConn == nil {
		log.Fatalf("Error conectando a RabbitMQ: nil connection")
	}
	defer rabbitConn.Close()

	// Inicializar repositorios
	searchRepo := repositories.NewSearchRepository(solrClient)

	// Inicializar servicios
	searchService := services.NewSearchService(searchRepo, cacheRepo)
	indexService := services.NewIndexService(searchRepo)

	// Inicializar event consumer
	eventConsumer, err := events.NewEventConsumer(rabbitConn, indexService)
	if err != nil {
		log.Fatalf("Error inicializando event consumer: %v", err)
	}

	// Iniciar consumidor de eventos en goroutine
	go func() {
		log.Println("🎧 Iniciando consumidor de eventos...")
		if err := eventConsumer.Start(); err != nil {
			log.Printf("Error en event consumer: %v", err)
		}
	}()

	// Inicializar controladores
	searchController := controllers.NewSearchController(searchService)

	// Provide a lightweight local index controller (stub) so the binary compiles even if
	// controllers.NewIndexController is not implemented; replace with real implementation later.
	indexController := struct {
		IndexRoom       func(*gin.Context)
		FullReindex     func(*gin.Context)
		DeleteFromIndex func(*gin.Context)
	}{
		IndexRoom: func(c *gin.Context) {
			id := c.Param("id")
			// TODO: call indexService to index a single room by id
			c.JSON(501, gin.H{"error": "IndexRoom not implemented", "id": id})
		},
		FullReindex: func(c *gin.Context) {
			// TODO: call indexService to perform full reindex
			c.JSON(501, gin.H{"error": "FullReindex not implemented"})
		},
		DeleteFromIndex: func(c *gin.Context) {
			id := c.Param("id")
			// TODO: call indexService to delete room from index by id
			c.JSON(501, gin.H{"error": "DeleteFromIndex not implemented", "id": id})
		},
	}

	// Configurar router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "search-api",
		})
	})

	// API Routes
	api := router.Group("/api/v1")
	{
		// Búsqueda de habitaciones
		search := api.Group("/search")
		{
			search.GET("/rooms", searchController.SearchRooms)
			search.GET("/rooms/suggestions", searchController.GetSuggestions)
			search.GET("/rooms/facets", searchController.GetFacets)
		}

		// Indexación manual (admin)
		admin := api.Group("/admin")
		{
			admin.POST("/index/room/:id", indexController.IndexRoom)
			admin.POST("/index/rooms/full", indexController.FullReindex)
			admin.DELETE("/index/room/:id", indexController.DeleteFromIndex)
		}
	}

	// Iniciar servidor
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Search API iniciando en puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
