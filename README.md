#  Distributed Rate-Limited API Gateway

A high-performance, cloud-native API Gateway designed to handle authentication, distributed rate limiting, and request proxying with low latency.

##  Architecture Features
- **Language:** Go (Golang) for high-concurrency processing.
- **Communication:** gRPC for internal service-to-service communication.
- **Rate Limiting:** Distributed Sliding Window algorithm using Redis & Lua scripts.
- **Storage:** PostgreSQL for persistent user data & API key management.
- **Observability:** Structured logging (Zap) and Prometheus metrics.
- **Resiliency:** Circuit Breaker and Graceful Shutdown.

## Tech Stack
- **Backend:** Go (Standard Library + gRPC)
- **Cache/In-memory:** Redis
- **Database:** PostgreSQL
- **DevOps:** Docker, Docker-compose, (Upcoming: Kubernetes)

---
