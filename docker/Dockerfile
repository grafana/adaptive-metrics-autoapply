FROM golang:1.27.2-bookworm@sha256:5cf287a799e6b94384bad13d16b14904c531f51ba65792237e122ce42b392f61 AS build

COPY . /opt
WORKDIR /opt

RUN CGO_ENABLED=0 GOOS=linux go build -o /adaptive-metrics-auto-apply ./cmd/adaptive-metrics/

FROM gcr.io/distroless/base-debian11@sha256:ac69aa622ea5dcbca0803ca877d47d069f51bd4282d5c96977e0390d7d256455

WORKDIR /

COPY --from=build /adaptive-metrics-auto-apply /adaptive-metrics-auto-apply

ENTRYPOINT ["/adaptive-metrics-auto-apply"]
