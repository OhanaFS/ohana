FROM golang:1.18-alpine AS builder

RUN set -ex; \
    mkdir -p /src; \
    apk add --update --no-cache \
    build-base git wget ca-certificates tzdata nodejs npm openssh; \
    rm -rf /var/cache/apk/*; \
    update-ca-certificates; \
    npm install -g yarn; \
    node --version;

WORKDIR /src
COPY go.mod go.sum Makefile ./
COPY web/package.json web/yarn.lock ./web/
RUN make deps
COPY . .
RUN make

FROM golang:1.18-alpine
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=1001:1001 /src/bin/ohana /app/ohana
CMD ["/app/ohana"]
