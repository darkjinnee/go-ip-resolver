package resolver

import (
	"context"
	"net"
	"time"
)

type Result struct {
	Domain string   `json:"domain"`
	IPs    []string `json:"ips"`
}

type Resolver struct {
	timeout time.Duration
}

func New(timeout time.Duration) *Resolver {
	return &Resolver{
		timeout: timeout,
	}
}

func (r *Resolver) ResolveDomain(domain string) []string {
	return r.ResolveDomainWithFilter(domain, "")
}

func (r *Resolver) ResolveDomainWithFilter(domain string, ipType string) []string {
	ips := []string{}

	// Создаем resolver с таймаутом
	resolver := &net.Resolver{
		PreferGo: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	addrs, err := resolver.LookupIPAddr(ctx, domain)
	if err == nil {
		for _, ip := range addrs {
			ipStr := ip.IP.String()

			// Фильтруем по типу IP
			switch ipType {
			case "ipv4":
				if ip.IP.To4() != nil {
					ips = append(ips, ipStr)
				}
			case "ipv6":
				if ip.IP.To4() == nil {
					ips = append(ips, ipStr)
				}
			default:
				ips = append(ips, ipStr)
			}
		}
	}

	return ips
}

func (r *Resolver) ResolveDomains(domains []string) []Result {
	return r.ResolveDomainsWithFilter(domains, "")
}

func (r *Resolver) ResolveDomainsWithFilter(domains []string, ipType string) []Result {
	results := make([]Result, 0, len(domains))

	for _, domain := range domains {
		ips := r.ResolveDomainWithFilter(domain, ipType)
		results = append(results, Result{
			Domain: domain,
			IPs:    ips,
		})
	}

	return results
}
