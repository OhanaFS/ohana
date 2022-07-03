FROM golang:1.18-alpine AS builder

RUN set -ex; \
    mkdir -p /src; \
    apk add --update --no-cache build-base git upx ca-certificates tzdata; \
    rm -rf /var/cache/apk/*; \
    update-ca-certificates;
WORKDIR /src
COPY . .
RUN make

FROM scratch
USER 1001
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=1001:1001 /src/bin/ohana /app/ohana
CMD ["/app/ohana"]
