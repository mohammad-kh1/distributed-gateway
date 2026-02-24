package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"github.com/mohammad-kh1/distributed-gateway/internal/gateway"
	"github.com/mohammad-kh1/distributed-gateway/internal/ratelimit"
)

func main() {
	// connect to Auth Service(grpc client)
	authClient := gateway.NewAuthClient("localhost:50051")

	authCache := gateway.NewAuthCache()
	limiter := ratelimit.NewRedisLimiter("localhost:6379")

	// config Circuit Breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:"Auth-Service",
		MaxRequest: 3,
		Interval: 5 * time.Second,
		Timeout: 10 * time.Second // how many Circuit will be open
	})


	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// 1 check for cache
		res , found := authCache.Get(token)
		if !found{
			// 2. if not found connect to Circuit Breaker
			body , err := cb.Execute(func()(interface{},error){
				return authClient.Authenticate(token)
			})

			if err != nil {
				http.Error(w , "Service Temporarily Unavilable (CB OPEN" , http.StatusServiceUnavailable)
				return
			}
			res = body.(*proto.VerifyResponse)
			// save in cache for future requests
			authCache.Set(token , res)
		}

		if !res.Authorized{
			http.Error(w , "َUnauthorized" ,http.StatusUnauthorized)
			return
		}



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
