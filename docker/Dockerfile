FROM golang:1.22.5-bookworm AS build

COPY . /opt
WORKDIR /opt

RUN CGO_ENABLED=0 GOOS=linux go build -o /adaptive-metrics-auto-apply ./cmd/adaptive-metrics/

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /adaptive-metrics-auto-apply /adaptive-metrics-auto-apply

ENTRYPOINT ["/adaptive-metrics-auto-apply"]
