FROM hashicorp/terraform:1.7

RUN apk add --no-cache jq

COPY --chmod=755 apply_recommendations.sh /opt/apply_recommendations.sh

ENTRYPOINT ["/opt/apply_recommendations.sh"]
