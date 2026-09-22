FROM golang:1.27.1-bookworm@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195 AS build

COPY . /opt
WORKDIR /opt

RUN CGO_ENABLED=0 GOOS=linux go build -o /adaptive-metrics-auto-apply ./cmd/adaptive-metrics/

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /adaptive-metrics-auto-apply /adaptive-metrics-auto-apply

ENTRYPOINT ["/adaptive-metrics-auto-apply"]
