ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG GO_BUILD_ARGS=""

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ${GO_BUILD_ARGS} -o /build/frontend-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:cba10d7abd3e203428e86f5b2d7fd5eb7d8987c387864ae4996cf97191b33764 AS runtime
WORKDIR /service
COPY --from=builder /build/frontend-service ./service
COPY internal/frontend/templates /service/frontend/internal/frontend/templates
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8081
ENTRYPOINT ["./service", "-config", "/config/config.yaml"]
