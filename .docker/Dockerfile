FROM golang:1.18-alpine AS builder

RUN set -ex; \
    mkdir -p /src; \
    apk add --update --no-cache build-base git wget ca-certificates tzdata; \
    rm -rf /var/cache/apk/*; \
    update-ca-certificates; \
    wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz; \
    tar -xJvf upx-3.96-amd64_linux.tar.xz -C /usr/local/bin; \
    mv /usr/local/bin/upx-3.96-amd64_linux/upx /usr/local/bin/upx; \
    rm -rvf upx-3.96-amd64_linux.tar.xz /usr/local/bin/upx-3.96-amd64_linux;
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make

FROM scratch
USER 1001
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=1001:1001 /src/bin/ohana /app/ohana
CMD ["/app/ohana"]
