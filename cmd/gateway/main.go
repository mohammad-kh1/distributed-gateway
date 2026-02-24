package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mohammad-kh1/distributed-gateway/internal/gateway"
	"github.com/mohammad-kh1/distributed-gateway/internal/ratelimit"
)

func main() {
	// connect to Auth Service(grpc client)
	authClient := gateway.NewAuthClient("localhost:50051")
	limiter := ratelimit.NewRedisLimiter("localhost:6379")

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		res, err := authClient.Authenticate(token)
		if err != nil || !res.Authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		allowed, err := limiter.IsAllowed(r.Context(), res.UserId, int(res.RateLimit))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Too Many Requests! Slow down.", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Welcome User: %s! Your rate limit is %d requests per minute.\n",
			res.UserId,
			res.RateLimit,
		)
	})

	log.Println("Gateway is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
