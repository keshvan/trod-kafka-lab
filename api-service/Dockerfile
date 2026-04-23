FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o api-service ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/api-service .
EXPOSE 8080
ENTRYPOINT ["./api-service"]