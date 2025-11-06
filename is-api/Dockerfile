# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY is-api/go.mod is-api/go.sum* ./is-api/
COPY is-worker/go.mod is-worker/go.sum* ./is-worker/

# Download all dependencies for the workspace
RUN go mod download

# Copy all source code
COPY pkg/ ./pkg/
COPY is-api/ ./is-api/
COPY is-worker/ ./is-worker/

# Build the API
WORKDIR /app/is-api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /app/bin/api .
EXPOSE 8000

CMD ["./api"]