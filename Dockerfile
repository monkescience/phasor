ARG VERSION

# Builder stage - builds both backend and frontend
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG GO_BUILD_ARGS=""

WORKDIR /build

# Copy go workspace and module files
COPY ./go.work ./go.work.sum ./
COPY ./backend/go.mod ./backend/go.sum ./backend/
COPY ./frontend/go.mod ./frontend/go.sum ./frontend/

# Download dependencies for both services
RUN cd backend && go mod download
RUN cd frontend && go mod download

# Copy source code
COPY ./backend ./backend
COPY ./frontend ./frontend

# Build backend service
RUN cd backend && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o /build/backend-service ./cmd/main.go

# Build frontend service
RUN cd frontend && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o /build/frontend-service ./cmd/main.go

# Backend runtime stage
FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412 AS backend-runtime
WORKDIR /service
COPY --from=builder /build/backend-service ./service
COPY ./backend/config/config.yaml /config/config.yaml
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8080
CMD ["./service", "-config", "/config/config.yaml"]

# Frontend runtime stage
FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412 AS frontend-runtime
WORKDIR /service
COPY --from=builder /build/frontend-service ./service
COPY ./frontend/config/config.docker.yaml /config/config.yaml
COPY ./frontend/internal/frontend/templates /service/frontend/internal/frontend/templates
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8081
CMD ["./service", "-config", "/config/config.yaml"]
