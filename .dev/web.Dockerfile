FROM node:16-alpine

ARG MY_UID
ARG MY_GID

# Create a user
RUN adduser \
        -u $MY_UID \
        -g '' \
        --disabled-password \
        user || true;
USER $MY_UID

COPY --chown=$MY_UID ./.dev/web.run.sh /run.sh
RUN chmod +x /run.sh

WORKDIR /home/user/app/web
CMD ["/run.sh"]
