ARG VERSION

# Builder stage
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG GO_BUILD_ARGS=""

WORKDIR /build

# Copy module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build backend service
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o /build/backend-service ./cmd/main.go

# Runtime stage
FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412 AS runtime
WORKDIR /service

# Copy binary from builder
COPY --from=builder /build/backend-service ./service

# Set version environment variable
ARG VERSION
ENV VERSION=${VERSION}

# Expose port
EXPOSE 8080

# Run service
CMD ["./service", "-config", "/config/config.yaml"]
