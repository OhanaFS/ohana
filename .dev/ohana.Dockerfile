FROM golang:1.18-alpine

ARG MY_UID
ARG MY_GID

RUN set -ex; \
    mkdir -p /src; \
    apk add --update --no-cache \
        build-base git wget ca-certificates tzdata nodejs npm openssh; \
    rm -rf /var/cache/apk/*; \
    update-ca-certificates; \
    npm install -g yarn; \
    node --version;

# Create a user
RUN adduser \
        -u $MY_UID \
        -g '' \
        --disabled-password \
        user || true; \
    chown -R $MY_UID:$MY_GID /src;
USER $MY_UID

# Install gin for hot reloading
RUN go install github.com/codegangsta/gin@latest

# Set up the private module
RUN mkdir -pm 0700 /home/user/.ssh
ENV GOPRIVATE=github.com/OhanaFS/stitch
COPY --chown=user ./.docker/deploy.key /home/user/.ssh/id_ed25519
RUN set -ex; \
    chmod 600 ~/.ssh/id_ed25519; \
    ssh-keyscan github.com >> ~/.ssh/known_hosts; \
    git config --global \
        url."ssh://git@github.com/OhanaFS/stitch".insteadOf \
            "https://github.com/OhanaFS/stitch";

# Pre-cache project dependencies
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download -x

COPY --chown=$MY_UID ./.dev/ohana.run.sh /run.sh
RUN chmod +x /run.sh

WORKDIR /home/user/app
CMD ["/run.sh"]
