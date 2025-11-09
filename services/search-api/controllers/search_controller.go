package controllers

import (
	"net/http"
	"search-api/domain"
	"search-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchService *services.SearchService
}

func NewSearchController(searchService *services.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

// SearchRooms maneja búsquedas de habitaciones
// @Summary Buscar habitaciones
// @Description Busca habitaciones con filtros avanzados
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Búsqueda por texto"
// @Param min_price query number false "Precio mínimo"
// @Param max_price query number false "Precio máximo"
// @Param min_capacity query int false "Capacidad mínima"
// @Param room_type query string false "Tipo de habitación"
// @Param is_available query bool false "Solo disponibles"
// @Param floor query int false "Piso"
// @Param has_wifi query bool false "Tiene WiFi"
// @Param has_ac query bool false "Tiene aire acondicionado"
// @Param has_tv query bool false "Tiene TV"
// @Param has_minibar query bool false "Tiene minibar"
// @Param page query int false "Número de página (default: 1)"
// @Param limit query int false "Resultados por página (default: 10, max: 100)"
// @Param sort query string false "Ordenamiento (price_asc, price_desc, capacity_asc, capacity_desc)"
// @Success 200 {object} domain.SearchResults
// @Router /api/v1/search/rooms [get]
func (c *SearchController) SearchRooms(ctx *gin.Context) {
	// Parsear parámetros
	params := domain.SearchParams{
		Query:       ctx.Query("q"),
		MinPrice:    parseFloat(ctx.Query("min_price"), 0),
		MaxPrice:    parseFloat(ctx.Query("max_price"), 0),
		MinCapacity: parseInt(ctx.Query("min_capacity"), 0),
		RoomType:    ctx.Query("room_type"),
		Status:      ctx.Query("status"),
		IsAvailable: parseBool(ctx.Query("is_available"), false),
		Floor:       parseInt(ctx.Query("floor"), 0),
		HasWifi:     parseBool(ctx.Query("has_wifi"), false),
		HasAC:       parseBool(ctx.Query("has_ac"), false),
		HasTV:       parseBool(ctx.Query("has_tv"), false),
		HasMinibar:  parseBool(ctx.Query("has_minibar"), false),
	}

	// Paginación
	page := parseInt(ctx.Query("page"), 1)
	limit := parseInt(ctx.Query("limit"), 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	params.Start = (page - 1) * limit
	params.Rows = limit

	// Ordenamiento
	sortParam := ctx.Query("sort")
	switch sortParam {
	case "price_asc":
		params.Sort = "price_per_night asc"
	case "price_desc":
		params.Sort = "price_per_night desc"
	case "capacity_asc":
		params.Sort = "capacity asc"
	case "capacity_desc":
		params.Sort = "capacity desc"
	default:
		params.Sort = "room_number asc"
	}

	// Facetas (solo si se solicitan explícitamente)
	params.IncludeFacets = parseBool(ctx.Query("include_facets"), false)

	// Buscar
	results, err := c.searchService.SearchRooms(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error realizando búsqueda",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

// GetSuggestions obtiene sugerencias de autocompletado
// @Summary Obtener sugerencias
// @Description Obtiene sugerencias para autocompletado basadas en un prefijo
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Prefijo de búsqueda"
// @Param limit query int false "Límite de sugerencias (default: 10, max: 20)"
// @Success 200 {object} domain.SuggestionResponse
// @Router /api/v1/search/rooms/suggestions [get]
func (c *SearchController) GetSuggestions(ctx *gin.Context) {
	prefix := ctx.Query("q")
	if prefix == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'q' es requerido",
		})
		return
	}

	limit := parseInt(ctx.Query("limit"), 10)
	if limit > 20 {
		limit = 20
	}

	suggestions, err := c.searchService.GetSuggestions(prefix, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error obteniendo sugerencias",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.SuggestionResponse{
		Suggestions: suggestions,
	})
}

// GetFacets obtiene facetas para filtros dinámicos
// @Summary Obtener facetas
// @Description Obtiene facetas para construir filtros dinámicos
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {object} domain.FacetResponse
// @Router /api/v1/search/rooms/facets [get]
func (c *SearchController) GetFacets(ctx *gin.Context) {
	facets, err := c.searchService.GetFacets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error obteniendo facetas",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, facets)
}

// Helper functions
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

func parseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

func parseBool(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return val
}
