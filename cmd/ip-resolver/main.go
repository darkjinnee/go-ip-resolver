package main

import (
	"flag"
	"log"
	"time"

	"go-ip-resolver/internal/cache"
	"go-ip-resolver/internal/config"
	"go-ip-resolver/internal/resolver"
	"go-ip-resolver/internal/transport"
)

func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	var (
		configPath     = flag.String("config", "configs/domains.json", "path to domains config file")
		addr           = flag.String("addr", ":8080", "server address")
		timeout        = flag.Duration("timeout", 2*time.Second, "DNS resolution timeout")
		cacheTTL       = flag.Duration("cache-ttl", 15*time.Minute, "cache TTL")
		updateInterval = flag.Duration("update-interval", 10*time.Minute, "cache update interval")
	)
	flag.Parse()

	// Загружаем конфигурацию
	log.Printf("Loading config from: %s", *configPath)
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded %d groups: %v", len(*cfg), getKeys(*cfg))

	// Создаем резолвер
	res := resolver.New(*timeout)

	// Создаем кэш
	cache := cache.New(res, *cacheTTL)
	log.Printf("Cache initialized with TTL: %v", *cacheTTL)

	// Запускаем периодическое обновление кэша
	go cache.StartPeriodicUpdate(*cfg, *updateInterval)
	log.Printf("Started periodic cache update every: %v", *updateInterval)

	// Создаем сервер
	server := transport.New(cfg, res, cache)

	// Запускаем сервер
	log.Printf("Starting server on %s", *addr)
	if err := server.Start(*addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
