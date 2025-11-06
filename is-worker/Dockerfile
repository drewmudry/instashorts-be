# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY is-worker/go.mod is-worker/go.sum* ./is-worker/
COPY is-api/go.mod is-api/go.sum* ./is-api/ 

# Download all dependencies for the workspace
RUN go mod download

# Copy all source code
COPY pkg/ ./pkg/
COPY is-worker/ ./is-worker/
COPY is-api/ ./is-api/ 

# Build the worker
WORKDIR /app/is-worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /app/bin/worker .
EXPOSE 5100

CMD ["./worker"]