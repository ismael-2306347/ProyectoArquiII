package config

import (
	"fmt"
	"net/http"
	"time"
)

type SolrClient struct {
	BaseURL string
	Client  *http.Client
}

func NewSolrClient(baseURL string) *SolrClient {
	return &SolrClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SolrClient) HealthCheck() error {
	url := fmt.Sprintf("%s/admin/ping", s.BaseURL)
	resp, err := s.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
