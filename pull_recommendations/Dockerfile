FROM alpine:3.14

RUN apk add --no-cache jq curl

COPY --chmod=755 pull_recommendations.sh /opt/pull_recommendations.sh

ENTRYPOINT ["/opt/pull_recommendations.sh"]
