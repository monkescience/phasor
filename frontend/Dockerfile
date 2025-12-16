ARG VERSION

FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder
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

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /service
COPY --from=builder /build/frontend-service ./service
COPY internal/frontend/templates /service/frontend/internal/frontend/templates
ARG VERSION
ENV VERSION=${VERSION}
EXPOSE 8081
ENTRYPOINT ["./service", "-config", "/config/config.yaml"]
