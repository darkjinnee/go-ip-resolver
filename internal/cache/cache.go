package cache

import (
	"sync"
	"time"

	"go-ip-resolver/internal/resolver"
)

type CacheEntry struct {
	Data      []resolver.Result
	ExpiresAt time.Time
}

type Cache struct {
	mu       sync.RWMutex
	entries  map[string]map[string]CacheEntry // [group][ipType] -> CacheEntry
	ttl      time.Duration
	resolver *resolver.Resolver
	stopCh   chan struct{}
}

func New(resolver *resolver.Resolver, ttl time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]map[string]CacheEntry),
		ttl:      ttl,
		resolver: resolver,
		stopCh:   make(chan struct{}),
	}

	// Запускаем горутину для очистки устаревших записей
	go c.cleanup()

	return c
}

func (c *Cache) Get(group, ipType string) ([]resolver.Result, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	groupCache, exists := c.entries[group]
	if !exists {
		return nil, false
	}

	entry, exists := groupCache[ipType]
	if !exists {
		return nil, false
	}

	// Проверяем, не истек ли TTL
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

func (c *Cache) Set(group, ipType string, data []resolver.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.entries[group] == nil {
		c.entries[group] = make(map[string]CacheEntry)
	}

	c.entries[group][ipType] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) UpdateGroup(group string, domains []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Обновляем для всех типов IP
	ipTypes := []string{"", "ipv4", "ipv6"}

	for _, ipType := range ipTypes {
		data := c.resolver.ResolveDomainsWithFilter(domains, ipType)

		if c.entries[group] == nil {
			c.entries[group] = make(map[string]CacheEntry)
		}

		c.entries[group][ipType] = CacheEntry{
			Data:      data,
			ExpiresAt: time.Now().Add(c.ttl),
		}
	}
}

func (c *Cache) UpdateAllGroups(groups map[string][]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Обновляем все группы
	for group, domains := range groups {
		ipTypes := []string{"", "ipv4", "ipv6"}

		if c.entries[group] == nil {
			c.entries[group] = make(map[string]CacheEntry)
		}

		for _, ipType := range ipTypes {
			data := c.resolver.ResolveDomainsWithFilter(domains, ipType)

			c.entries[group][ipType] = CacheEntry{
				Data:      data,
				ExpiresAt: time.Now().Add(c.ttl),
			}
		}
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute) // Очистка каждые 5 минут
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for group, groupCache := range c.entries {
		for ipType, entry := range groupCache {
			if now.After(entry.ExpiresAt) {
				delete(groupCache, ipType)
			}
		}

		// Удаляем пустые группы
		if len(groupCache) == 0 {
			delete(c.entries, group)
		}
	}
}

func (c *Cache) StartPeriodicUpdate(groups map[string][]string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Первоначальное заполнение кэша
	c.UpdateAllGroups(groups)

	for {
		select {
		case <-ticker.C:
			c.UpdateAllGroups(groups)
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.stopCh)
}

func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := map[string]interface{}{
		"groups_count":  len(c.entries),
		"total_entries": 0,
		"groups":        make(map[string]int),
	}

	for group, groupCache := range c.entries {
		stats["groups"].(map[string]int)[group] = len(groupCache)
		stats["total_entries"] = stats["total_entries"].(int) + len(groupCache)
	}

	return stats
}
