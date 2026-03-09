# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/configs ./configs

EXPOSE 8080

CMD ["./server"]
