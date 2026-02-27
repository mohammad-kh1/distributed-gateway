FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/auth/main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /auth-service /auth-service
EXPOSE 50051
CMD ["/auth-service"]
