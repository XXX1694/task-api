package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type client struct {
	lastSeen time.Time
	count    int
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

const (
	rateLimit  = 10
	ratePeriod = time.Minute
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	go cleanupClients()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		// Убираем порт из адреса
		if idx := len(ip) - 1; idx >= 0 {
			for i := idx; i >= 0; i-- {
				if ip[i] == ':' {
					ip = ip[:i]
					break
				}
			}
		}

		mu.Lock()

		c, exists := clients[ip]
		if !exists {
			clients[ip] = &client{
				lastSeen: time.Now(),
				count:    1,
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(c.lastSeen) > ratePeriod {
			c.count = 1
			c.lastSeen = time.Now()
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if c.count >= rateLimit {
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}

		c.count++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func cleanupClients() {
	for {
		time.Sleep(ratePeriod)

		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > ratePeriod {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
