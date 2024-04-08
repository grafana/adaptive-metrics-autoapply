FROM hashicorp/terraform:1.7

RUN apk add --no-cache jq

COPY --chmod 755 autoapply.sh /opt/autoapply.sh

ENTRYPOINT ['/opt/autoapply.sh']
