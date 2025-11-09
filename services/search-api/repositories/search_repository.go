package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"search-api/config"
	"search-api/domain"
	"strings"
)

type SearchRepository struct {
	solrURL string
	client  *http.Client
}

func NewSearchRepository(solrClient *config.SolrClient) *SearchRepository {
	return &SearchRepository{
		solrURL: solrClient.BaseURL,
		client:  solrClient.Client,
	}
}

// IndexRoom indexa una habitación en Solr
func (r *SearchRepository) IndexRoom(room *domain.RoomDocument) error {
	url := fmt.Sprintf("%s/update?commit=true", r.solrURL)

	data, err := json.Marshal([]interface{}{room})
	if err != nil {
		return fmt.Errorf("error marshaling room: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error indexing room: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteRoom elimina una habitación del índice
func (r *SearchRepository) DeleteRoom(roomID string) error {
	url := fmt.Sprintf("%s/update?commit=true", r.solrURL)

	deleteQuery := map[string]interface{}{
		"delete": map[string]string{
			"id": roomID,
		},
	}

	data, err := json.Marshal(deleteQuery)
	if err != nil {
		return fmt.Errorf("error marshaling delete query: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting room: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchRooms busca habitaciones con filtros
func (r *SearchRepository) SearchRooms(params domain.SearchParams) (*domain.SearchResults, error) {
	url := fmt.Sprintf("%s/select", r.solrURL)

	// Construir query
	query := r.buildQuery(params)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error searching rooms: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("solr returned status %d: %s", resp.StatusCode, string(body))
	}

	var solrResp struct {
		Response struct {
			NumFound int                   `json:"numFound"`
			Docs     []domain.RoomDocument `json:"docs"`
		} `json:"response"`
		Facets map[string]interface{} `json:"facet_counts,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &domain.SearchResults{
		TotalResults: solrResp.Response.NumFound,
		Results:      solrResp.Response.Docs,
		Facets:       solrResp.Facets,
	}, nil
}

// buildQuery construye los parámetros de búsqueda de Solr
func (r *SearchRepository) buildQuery(params domain.SearchParams) map[string]string {
	query := make(map[string]string)

	// Query principal
	if params.Query != "" {
		query["q"] = params.Query
	} else {
		query["q"] = "*:*"
	}

	// Filtros
	var filters []string

	if params.MinPrice > 0 {
		filters = append(filters, fmt.Sprintf("price_per_night:[%f TO *]", params.MinPrice))
	}

	if params.MaxPrice > 0 {
		filters = append(filters, fmt.Sprintf("price_per_night:[* TO %f]", params.MaxPrice))
	}

	if params.MinCapacity > 0 {
		filters = append(filters, fmt.Sprintf("capacity:[%d TO *]", params.MinCapacity))
	}

	if params.RoomType != "" {
		filters = append(filters, fmt.Sprintf("room_type:\"%s\"", params.RoomType))
	}

	if params.Status != "" {
		filters = append(filters, fmt.Sprintf("status:\"%s\"", params.Status))
	}

	if params.IsAvailable {
		filters = append(filters, "is_available:true")
	}

	if len(filters) > 0 {
		query["fq"] = strings.Join(filters, " AND ")
	}

	// Paginación
	query["start"] = fmt.Sprintf("%d", params.Start)
	query["rows"] = fmt.Sprintf("%d", params.Rows)

	// Ordenamiento
	if params.Sort != "" {
		query["sort"] = params.Sort
	}

	// Formato de respuesta
	query["wt"] = "json"

	// Facetas (para filtros dinámicos)
	if params.IncludeFacets {
		query["facet"] = "true"
		query["facet.field"] = "room_type"
		query["facet.field"] = "status"
		query["facet.range"] = "price_per_night"
		query["f.price_per_night.facet.range.start"] = "0"
		query["f.price_per_night.facet.range.end"] = "1000"
		query["f.price_per_night.facet.range.gap"] = "100"
	}

	return query
}

// GetSuggestions obtiene sugerencias de autocompletado
func (r *SearchRepository) GetSuggestions(prefix string, limit int) ([]string, error) {
	url := fmt.Sprintf("%s/select", r.solrURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("room_number:%s* OR room_type:%s*", prefix, prefix))
	q.Add("rows", fmt.Sprintf("%d", limit))
	q.Add("fl", "room_number,room_type")
	q.Add("wt", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting suggestions: %w", err)
	}
	defer resp.Body.Close()

	var solrResp struct {
		Response struct {
			Docs []domain.RoomDocument `json:"docs"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	suggestions := make([]string, 0)
	seen := make(map[string]bool)

	for _, doc := range solrResp.Response.Docs {
		if !seen[doc.RoomNumber] {
			suggestions = append(suggestions, doc.RoomNumber)
			seen[doc.RoomNumber] = true
		}
		if !seen[doc.RoomType] {
			suggestions = append(suggestions, doc.RoomType)
			seen[doc.RoomType] = true
		}
	}

	return suggestions, nil
}
