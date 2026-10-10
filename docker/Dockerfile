FROM golang:1.27.1-bookworm@sha256:8d48e12ec56735e9358640898b9d9b9fcca110612ed8a5567438c0a1baa24e66 AS build

COPY . /opt
WORKDIR /opt

RUN CGO_ENABLED=0 GOOS=linux go build -o /adaptive-metrics-auto-apply ./cmd/adaptive-metrics/

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /adaptive-metrics-auto-apply /adaptive-metrics-auto-apply

ENTRYPOINT ["/adaptive-metrics-auto-apply"]
