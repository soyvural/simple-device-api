package store

import (
	"sync"

	"github.com/soyvural/simple-device-api/types"
)

const limit = 10000

type cache struct {
	mu   sync.RWMutex
	data map[string]types.Device
}

func NewCache() *cache {
	return &cache{
		data: make(map[string]types.Device),
	}
}

func (c *cache) Get(id string) *types.Device {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if d, ok := c.data[id]; ok {
		return &types.Device{
			ID:    d.ID,
			Name:  d.Name,
			Brand: d.Brand,
			Model: d.Model,
		}
	}
	return nil
}

func (c *cache) Insert(d types.Device) *types.Device {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.data) >= limit {
		return nil
	}
	c.data[d.ID] = d
	return &types.Device{
		ID:    d.ID,
		Name:  d.Name,
		Brand: d.Brand,
		Model: d.Model,
	}
}

func (c *cache) Delete(id string) *types.Device {
	c.mu.Lock()
	defer c.mu.Unlock()
	if d, ok := c.data[id]; ok {
		delete(c.data, id)
		return &types.Device{
			ID:    d.ID,
			Name:  d.Name,
			Brand: d.Brand,
			Model: d.Model,
		}
	}
	return nil
}
