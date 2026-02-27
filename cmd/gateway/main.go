package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"github.com/mohammad-kh1/distributed-gateway/internal/gateway"
	"github.com/mohammad-kh1/distributed-gateway/internal/logger"
	"github.com/mohammad-kh1/distributed-gateway/internal/ratelimit"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Get().Sync()

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	authClient := gateway.NewAuthClient("auth-service:50051")
	authCache := gateway.NewAuthCache()
	limiter := ratelimit.NewRedisLimiter(redisAddr)

	cb := initCircuitBreaker()

	mainMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	mainMux.HandleFunc("/dashboard.html", handleDashboard)
	mainMux.HandleFunc("/admin/metrics", handleMetricsJSON)
	mainMux.HandleFunc("/admin/metrics/stream", handleMetricsStream)

	apiMux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		handleData(w, r, authClient, authCache, limiter, cb)
	})

	mainMux.Handle("/api/", gateway.LoggingMiddleware(apiMux))

	log.Println("Gateway is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mainMux))
}

func handleData(w http.ResponseWriter, r *http.Request, auth *gateway.AuthClient, cache *gateway.AuthCache, lim *ratelimit.RedisLimiter, cb *gobreaker.CircuitBreaker) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Authorization header required")
		return
	}

	res, found := cache.Get(token)
	if !found {
		result, err := cb.Execute(func() (interface{}, error) {
			return auth.Authenticate(token)
		})

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "Auth service unavailable (Circuit Breaker Open)")
			return
		}

		res = result.(*proto.VerifyResponse)
		cache.Set(token, res)
	}

	if !res.Authorized {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
		return
	}

	allowed, err := lim.IsAllowed(r.Context(), res.UserId, int(res.RateLimit))
	if err != nil || !allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, "Too Many Requests! Slow down.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome %s! Rate limit: %d/min\n", res.UserId, res.RateLimit)
}

func handleMetricsStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			return
		default:
			jsonData, _ := json.Marshal(gateway.GlobalStats)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	}
}

func handleMetricsJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(gateway.GlobalStats)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dashboard.html")
}

func initCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Auth-Service",
		ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 3 },
		Timeout:     10 * time.Second,
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Get().Info("Circuit breaker state change",
				zap.String("service", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
