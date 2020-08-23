# this works: FROM golang:1.14.6-buster as golang-builder
FROM golang:1.14.6-alpine3.12 as golang-builder

WORKDIR /go/src/github.com/micha37-martins/soundboard

# use --no-cache to not cache the index locally
RUN apk --no-cache add git build-base alsa-lib-dev alsa-utils alsa-utils-doc alsa-lib alsaconf

#RUN addgroup root audio

ADD go.mod .
ADD go.sum .
ADD Makefile .

RUN make deps

ADD cmd/ cmd/
ADD internal/ internal/
ADD soundfiles/ soundfiles/

RUN make build

ENTRYPOINT ["./soundboard"]
