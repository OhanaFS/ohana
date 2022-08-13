FROM golang:1.18-alpine AS builder

RUN set -ex; \
    mkdir -p /src; \
    apk add --update --no-cache \
        build-base git wget ca-certificates tzdata nodejs npm openssh; \
    rm -rf /var/cache/apk/*; \
    update-ca-certificates; \
    npm install -g yarn; \
    node --version;

# Set up the private module
ENV GOPRIVATE=github.com/OhanaFS/stitch
COPY ./.docker/deploy.key /root/.ssh/id_ed25519
RUN set -ex; \
    chmod 600 /root/.ssh/id_ed25519; \
    ssh-keyscan github.com >> /root/.ssh/known_hosts; \
    git config --global \
        url."ssh://git@github.com/OhanaFS/stitch".insteadOf \
            "https://github.com/OhanaFS/stitch";

WORKDIR /src
COPY go.mod go.sum Makefile ./
COPY web/package.json web/yarn.lock ./web/
RUN make deps
COPY . .
RUN make

FROM scratch
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=1001:1001 /src/bin/ohana /app/ohana
CMD ["/app/ohana"]
