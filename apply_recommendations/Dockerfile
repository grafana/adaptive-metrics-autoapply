FROM golang:1.22.5-bookworm AS build

COPY . /opt
WORKDIR /opt

RUN CGO_ENABLED=0 GOOS=linux go build -o /apply-recommendations cmd/apply/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /apply-recommendations /apply-recommendations

ENTRYPOINT ["/apply-recommendations"]
