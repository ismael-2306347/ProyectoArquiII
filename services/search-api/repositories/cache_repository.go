package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheRepository struct {
	client *memcache.Client
	ttl    time.Duration
}

func NewCacheRepository(host, port string, ttl time.Duration) *CacheRepository {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := memcache.New(addr)
	client.Timeout = 2 * time.Second
	client.MaxIdleConns = 10

	return &CacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *CacheRepository) Get(key string) (interface{}, error) {
	item, err := r.client.Get(key)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(item.Value, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CacheRepository) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(r.ttl.Seconds()),
	}

	return r.client.Set(item)
}

func (r *CacheRepository) Delete(key string) error {
	return r.client.Delete(key)
}

func (r *CacheRepository) Ping() error {
	return r.client.Ping()
}

// GetString obtiene un valor string del cache
func (r *CacheRepository) GetString(key string) (string, error) {
	item, err := r.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// SetString guarda un string en el cache
func (r *CacheRepository) SetString(key, value string) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(r.ttl.Seconds()),
	}
	return r.client.Set(item)
}
