# Distributed API Gateway with Go & gRPC

A distributed API Gateway built in Go. This system acts as a centralized entry point
---

## Key Features

* **gRPC Communication:** Leverages gRPC binary protocol for ultra-fast communication between the Gateway and internal Auth services.
* **Circuit Breaker Implementation:** Prevents cascading failures using the `gobreaker` library, ensuring system stability when internal services are down.
* **Distributed Rate Limiting:** User-based traffic control using **Redis** with a fixed-window algorithm to prevent API abuse.
* **Real-time Monitoring:** A live observability dashboard powered by **SSE (Server-Sent Events)** to monitor traffic and system health without page refreshes.
* **Multi-level Caching:** Reduces latency by implementing a local in-memory cache for authentication data at the Gateway level.
* **Atomic Metrics:** High-performance, lock-free metrics recording using Go's `sync/atomic` package for minimal overhead under heavy load.

---

## Tech Stack

* **Language:** Go (1.21+)
* **Communication:** gRPC & Protocol Buffers
* **Database:** PostgreSQL (User Persistence)
* **Cache & Rate Limiting:** Redis
* **Monitoring:** HTML5, Chart.js, SSE
* **Containerization:** Docker & Docker Compose
* **Logging:** Uber Zap (Structured Logging)

---

## High-Level Architecture

Incoming requests are captured by the Gateway and passed through a chain of middleware (Logging, Metrics, Authentication, and Rate Limiting) before reaching the final destination.



---

## Quick Start

### 1. Prerequisites
Ensure you have **Docker** and **Docker Compose** installed on your machine.

### 2. Clone and Run
```bash
# Clone the repository
git clone https://github.com/mohammad-kh1/distributed-gateway.git
cd distributed-gateway

# Spin up all services in containers
docker-compose up --build
