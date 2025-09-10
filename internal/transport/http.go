package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-ip-resolver/internal/cache"
	"go-ip-resolver/internal/config"
	"go-ip-resolver/internal/resolver"
)

type Server struct {
	cfg      *config.Config
	resolver *resolver.Resolver
	cache    *cache.Cache
}

func New(cfg *config.Config, resolver *resolver.Resolver, cache *cache.Cache) *Server {
	return &Server{
		cfg:      cfg,
		resolver: resolver,
		cache:    cache,
	}
}

func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		http.Error(w, "group parameter is required", http.StatusBadRequest)
		return
	}

	domains, exists := (*s.cfg)[group]
	if !exists {
		http.Error(w, fmt.Sprintf("group '%s' not found", group), http.StatusNotFound)
		return
	}

	// Получаем тип IP для фильтрации
	ipType := r.URL.Query().Get("type")
	if ipType != "" && ipType != "ipv4" && ipType != "ipv6" {
		http.Error(w, "type parameter must be 'ipv4' or 'ipv6'", http.StatusBadRequest)
		return
	}

	// Пытаемся получить данные из кэша
	results, found := s.cache.Get(group, ipType)
	if !found {
		// Если в кэше нет, резолвим и сохраняем
		results = s.resolver.ResolveDomainsWithFilter(domains, ipType)
		s.cache.Set(group, ipType, results)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleListGroups(w http.ResponseWriter, r *http.Request) {
	groups := make([]string, 0, len(*s.cfg))
	for group := range *s.cfg {
		groups = append(groups, group)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{"groups": groups}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleResolveAll(w http.ResponseWriter, r *http.Request) {
	// Получаем тип IP для фильтрации
	ipType := r.URL.Query().Get("type")
	if ipType != "" && ipType != "ipv4" && ipType != "ipv6" {
		http.Error(w, "type parameter must be 'ipv4' or 'ipv6'", http.StatusBadRequest)
		return
	}

	allResults := make(map[string][]resolver.Result)

	for group, domains := range *s.cfg {
		// Пытаемся получить данные из кэша
		results, found := s.cache.Get(group, ipType)
		if !found {
			// Если в кэше нет, резолвим и сохраняем
			results = s.resolver.ResolveDomainsWithFilter(domains, ipType)
			s.cache.Set(group, ipType, results)
		}
		allResults[group] = results
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allResults); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := s.cache.GetStats()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleResolveFlat(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		http.Error(w, "group parameter is required", http.StatusBadRequest)
		return
	}

	domains, exists := (*s.cfg)[group]
	if !exists {
		http.Error(w, fmt.Sprintf("group '%s' not found", group), http.StatusNotFound)
		return
	}

	// Получаем тип IP для фильтрации
	ipType := r.URL.Query().Get("type")
	if ipType != "" && ipType != "ipv4" && ipType != "ipv6" {
		http.Error(w, "type parameter must be 'ipv4' or 'ipv6'", http.StatusBadRequest)
		return
	}

	// Пытаемся получить данные из кэша
	results, found := s.cache.Get(group, ipType)
	if !found {
		// Если в кэше нет, резолвим и сохраняем
		results = s.resolver.ResolveDomainsWithFilter(domains, ipType)
		s.cache.Set(group, ipType, results)
	}

	// Собираем все IP в плоский список
	var allIPs []string
	for _, result := range results {
		allIPs = append(allIPs, result.IPs...)
	}

	// Устанавливаем Content-Type как text/plain
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Записываем каждый IP на отдельной строке
	for _, ip := range allIPs {
		fmt.Fprintln(w, ip)
	}
}

func (s *Server) handleResolveFlatAll(w http.ResponseWriter, r *http.Request) {
	// Получаем тип IP для фильтрации
	ipType := r.URL.Query().Get("type")
	if ipType != "" && ipType != "ipv4" && ipType != "ipv6" {
		http.Error(w, "type parameter must be 'ipv4' or 'ipv6'", http.StatusBadRequest)
		return
	}

	// Собираем все IP из всех групп в плоский список
	var allIPs []string
	for group, domains := range *s.cfg {
		// Пытаемся получить данные из кэша
		results, found := s.cache.Get(group, ipType)
		if !found {
			// Если в кэше нет, резолвим и сохраняем
			results = s.resolver.ResolveDomainsWithFilter(domains, ipType)
			s.cache.Set(group, ipType, results)
		}

		// Добавляем IP из текущей группы
		for _, result := range results {
			allIPs = append(allIPs, result.IPs...)
		}
	}

	// Устанавливаем Content-Type как text/plain
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Записываем каждый IP на отдельной строке
	for _, ip := range allIPs {
		fmt.Fprintln(w, ip)
	}
}

func (s *Server) handleResolveFlatAllIPv4(w http.ResponseWriter, r *http.Request) {
	// Собираем все IPv4 из всех групп в плоский список
	var allIPs []string
	for group, domains := range *s.cfg {
		// Пытаемся получить данные из кэша
		results, found := s.cache.Get(group, "ipv4")
		if !found {
			// Если в кэше нет, резолвим и сохраняем
			results = s.resolver.ResolveDomainsWithFilter(domains, "ipv4")
			s.cache.Set(group, "ipv4", results)
		}

		// Добавляем IP из текущей группы
		for _, result := range results {
			allIPs = append(allIPs, result.IPs...)
		}
	}

	// Устанавливаем Content-Type как text/plain
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Записываем каждый IP на отдельной строке
	for _, ip := range allIPs {
		fmt.Fprintln(w, ip)
	}
}

func (s *Server) handleResolveFlatAllIPv6(w http.ResponseWriter, r *http.Request) {
	// Собираем все IPv6 из всех групп в плоский список
	var allIPs []string
	for group, domains := range *s.cfg {
		// Пытаемся получить данные из кэша
		results, found := s.cache.Get(group, "ipv6")
		if !found {
			// Если в кэше нет, резолвим и сохраняем
			results = s.resolver.ResolveDomainsWithFilter(domains, "ipv6")
			s.cache.Set(group, "ipv6", results)
		}

		// Добавляем IP из текущей группы
		for _, result := range results {
			allIPs = append(allIPs, result.IPs...)
		}
	}

	// Устанавливаем Content-Type как text/plain
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Записываем каждый IP на отдельной строке
	for _, ip := range allIPs {
		fmt.Fprintln(w, ip)
	}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/resolve", s.handleResolve)
	mux.HandleFunc("/resolve-all", s.handleResolveAll)
	mux.HandleFunc("/resolve-flat", s.handleResolveFlat)
	mux.HandleFunc("/resolve-flat-all", s.handleResolveFlatAll)
	mux.HandleFunc("/resolve-flat-all-ipv4", s.handleResolveFlatAllIPv4)
	mux.HandleFunc("/resolve-flat-all-ipv6", s.handleResolveFlatAllIPv6)
	mux.HandleFunc("/groups", s.handleListGroups)
	mux.HandleFunc("/cache/stats", s.handleCacheStats)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server running on %s\n", addr)
	return server.ListenAndServe()
}
