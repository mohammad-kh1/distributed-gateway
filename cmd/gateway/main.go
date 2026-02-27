package main

import (
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	authClient := gateway.NewAuthClient("auth-service:50051")
	authCache := gateway.NewAuthCache()

	limiter := ratelimit.NewRedisLimiter(redisAddr)
	// Circuit Breaker configuration
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Auth-Service",
		ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 3 },
		Timeout:     10 * time.Second,
		Interval:    5 * time.Second,
		MaxRequests: 3,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Get().Info("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		handleData(w, r, authClient, authCache, limiter, cb)
	})

	loggingMux := gateway.LoggingMiddleware(mux)

	log.Println("Gateway is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", loggingMux))
}

// handleData now takes dependencies as parameters (cleaner, testable)
func handleData(
	w http.ResponseWriter,
	r *http.Request,
	authClient *gateway.AuthClient,
	authCache *gateway.AuthCache,
	limiter *ratelimit.RedisLimiter,
	cb *gobreaker.CircuitBreaker,
) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// 1. Check cache first (fast path)
	res, found := authCache.Get(token)
	if !found {
		// 2. Use circuit breaker to call auth service
		result, err := cb.Execute(func() (interface{}, error) {
			return authClient.Authenticate(token)
		})

		if err != nil {
			// Circuit open / timeout / too many failures → 503
			http.Error(w, "Auth service temporarily unavailable (circuit breaker open)", http.StatusServiceUnavailable)
			return
		}

		res = result.(*proto.VerifyResponse)

		authCache.Set(token, res)
	}

	// 3. Authorization check
	if !res.Authorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 4. Rate limiting
	allowed, err := limiter.IsAllowed(r.Context(), res.UserId, int(res.RateLimit))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Too Many Requests! Slow down.", http.StatusTooManyRequests)
		return
	}

	// 5. Success response
	fmt.Fprintf(w, "Welcome User: %s! Your rate limit is %d requests per minute.\n",
		res.UserId,
		res.RateLimit,
	)
}
