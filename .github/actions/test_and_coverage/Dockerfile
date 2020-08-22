# This Dockerfile adds required dependencies to the alpine image
FROM golang:1.14.6-alpine3.12 as golang-builder

# use --no-cache to not cache the index locally
RUN apk --no-cache add git build-base alsa-lib-dev alsa-utils alsa-utils-doc alsa-lib alsaconf

# Copies code file from action repository to the filesystem path `/` of the container
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
