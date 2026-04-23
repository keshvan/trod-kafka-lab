FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o data-service ./cmd/data
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/data-service ./data-service
COPY --from=builder /build/migrate ./migrate
COPY --from=builder /build/migrations ./migrations
EXPOSE 8080
CMD ["./data-service"]
