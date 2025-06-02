package cache

import (
	"digital-trainer/internal/models"
	"sync"
)

type Cache struct {
	metrics []models.Metrics
	mu      sync.Mutex
	maxSize int
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		metrics: make([]models.Metrics, 0, maxSize),
		maxSize: maxSize,
	}
}

func (c *Cache) AddMetric(metric models.Metrics) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.metrics) >= c.maxSize {
		c.metrics = c.metrics[1:]
	}
	c.metrics = append(c.metrics, metric)
}

func (c *Cache) GetMetrics() []models.Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.metrics
}
