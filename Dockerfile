# Stage 1: Build
# Uses golang:1.25.8-alpine as the builder image.
FROM golang:1.25.8-alpine AS builder

WORKDIR /build

# Download dependencies first (layer cache optimisation)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /build/paas-demo-app \
    ./cmd/server

# Stage 2: Runtime
# Uses distroless nonroot — no shell, no package manager, runs as UID 65532.
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the compiled binary from the builder stage
COPY --from=builder /build/paas-demo-app /app/paas-demo-app

# Expose the application port
EXPOSE 8080

# Run as nonroot (UID 65532) — satisfies restricted-v2 SCC
USER nonroot:nonroot

ENTRYPOINT ["/app/paas-demo-app"]
